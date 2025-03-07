package vinvc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/DIMO-Network/attestation-api/pkg/verifiable"
	pb "github.com/DIMO-Network/fetch-api/pkg/grpc"
	"github.com/DIMO-Network/model-garage/pkg/cloudevent"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var dincSource = common.HexToAddress("0x4F098Ea7cAd393365b4d251Dd109e791e6190239")

var notFound = errors.New("not found")

// FetchAPIService defines the interface Fetch API.
type FetchAPIService interface {
	// ListCloudEvents retrieves the most recent cloud events matching the provided search criteria.
	ListCloudEvents(ctx context.Context, filter *pb.SearchOptions, limit int32) ([]cloudevent.CloudEvent[json.RawMessage], error)
}

// VINVCService is a client for interacting with the VINVC service.
type VINVCService struct {
	fetchService     FetchAPIService
	logger           zerolog.Logger
	vinVCDataVersion string
	vehicleAddr      common.Address
	chainID          uint64
	trustedRecorders []string
}

// New creates a new instance of VINVCClient.
func New(fetchService FetchAPIService, settings *config.Settings, logger *zerolog.Logger) *VINVCService {
	return &VINVCService{
		fetchService:     fetchService,
		logger:           logger.With().Str("component", "vinvc-service").Logger(),
		vinVCDataVersion: settings.VINVCDataVersion,
		vehicleAddr:      settings.VehicleNFTAddress,
		chainID:          uint64(settings.DIMORegistryChainID),
		trustedRecorders: []string{
			cloudevent.EthrDID{
				ChainID:         uint64(settings.DIMORegistryChainID),
				ContractAddress: dincSource,
			}.String(),
		},
	}
}

// GetConfirmedVINVCs retrieves confirmed VINVCs for the provided vehicles.
// A VINVC is confirmed if:
// 1. It is not expired.
// 2. It was either recorded in the past week OR recorded by a trusted source.
// 3. The TokenId holds the latest recordedAt for the VIN.
func (v *VINVCService) GetConfirmedVINVCs(ctx context.Context, activeVehicles []*ch.Vehicle, weekNum int) (map[int64]struct{}, error) {
	validVinToTokenIDs, err := v.GetLatestValidVINVCs(ctx, activeVehicles, weekNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest valid VINVCs: %w", err)
	}
	confirmedVINs := v.ResolveVINConflicts(validVinToTokenIDs)

	return confirmedVINs, nil
}

// GetLatestValidVINVCs retrieves the latest VINVC for the provided vehicle tokenId in the given week.
func (v *VINVCService) GetLatestValidVINVCs(ctx context.Context, activeVehicles []*ch.Vehicle, weekNum int) (map[string][]*verifiable.VINSubject, error) {
	// map to track VINs and their associated subjects (only for valid VCs)
	validVinToCredSubjects := make(map[string][]*verifiable.VINSubject)

	// collect all valid VINVCs and their associated VINs
	for _, vehicle := range activeVehicles {
		logger := v.logger.With().Int64("vehicleTokenId", vehicle.TokenID).Logger()
		credSubject, err := v.getLatestValidVINVC(ctx, vehicle.TokenID, weekNum)
		if err != nil {
			if errors.Is(err, notFound) {
				logger.Warn().Msg("no VINVC found for vehicle")
				continue
			}
			return nil, fmt.Errorf("failed to get latest VINVC: %w", err)
		}

		vin := credSubject.VehicleIdentificationNumber
		validVinToCredSubjects[vin] = append(validVinToCredSubjects[vin], credSubject)
	}

	return validVinToCredSubjects, nil
}

// getLatestValidVINVC retrieves the latest VINVC for the provided vehicle tokenId in the given week.
// if a VINVC fails to decode, it is skipped and the next one is fetched.
func (v *VINVCService) getLatestValidVINVC(ctx context.Context, tokenId int64, weekNum int) (*verifiable.VINSubject, error) {
	// get time boundaries for the specified week
	endOfWeek := date.NumToWeekEnd(weekNum)
	startOfWeek := date.NumToWeekStart(weekNum)

	// initialize search time to the end of the week (plus a second for inclusive search)
	searchTime := endOfWeek.Add(time.Second)
	opts := &pb.SearchOptions{
		DataVersion: &wrapperspb.StringValue{Value: v.vinVCDataVersion},
		Type:        &wrapperspb.StringValue{Value: cloudevent.TypeVerifableCredential},
		Subject:     &wrapperspb.StringValue{Value: cloudevent.NFTDID{ChainID: v.chainID, ContractAddress: v.vehicleAddr, TokenID: uint32(tokenId)}.String()},
	}

	logger := v.logger.With().Int64("vehicleTokenId", tokenId).Logger()

	// continue searching until we reach the beginning of the week
	for startOfWeek.Before(searchTime) {
		// create search options with the current time boundary
		opts.Before = timestamppb.New(searchTime)

		// fetch the latest cloud event before the current search time
		cloudEvents, err := v.fetchService.ListCloudEvents(ctx, opts, 1)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, fmt.Errorf("no VINVC found for vehicle: %w", notFound)
			}
			return nil, fmt.Errorf("failed to get fetch latest VINVC: %w", err)
		}

		if len(cloudEvents) == 0 {
			// no more events found, exit the loop
			break
		}

		cloudEvent := cloudEvents[0]
		// update search time to look before the current event
		searchTime = cloudEvent.Time
		eventLogger := logger.With().Str("cloudEventId", cloudEvent.ID).Str("cloudEventSource", cloudEvent.Source).Logger()

		// parse and validate the credential
		var cred verifiable.Credential
		if err := json.Unmarshal(cloudEvent.Data, &cred); err != nil {
			eventLogger.Error().Err(err).Msg("failed to unmarshal VIN credential, skipping...")
			//
			continue
		}

		// parse the credential subject
		var credSubject verifiable.VINSubject
		if err := json.Unmarshal(cred.CredentialSubject, &credSubject); err != nil {
			eventLogger.Error().Err(err).Msg("failed to unmarshal VIN credential subject, skipping...")
			continue
		}

		// skip if VIN is empty
		if credSubject.VehicleIdentificationNumber == "" {
			eventLogger.Error().Msg("VINVC has empty VIN, skipping...")
			continue
		}

		// check if this is a valid VC for baseline issuance
		if !v.isValidVC(&logger, uint32(tokenId), &cred, &credSubject, weekNum) {
			eventLogger.Error().Any("credentialSubject", credSubject).Msg("VINVC did not meet validation criteria")
			continue
		}

		// found a valid credential, return it
		return &credSubject, nil
	}

	// no valid credential found within the week
	return nil, fmt.Errorf("no valid VINVC found for vehicle in week %d: %w", weekNum, notFound)
}

// ResolveVINConflicts resolves conflicts between VINs and their associated with multiple tokenIDs.
func (v *VINVCService) ResolveVINConflicts(vinToTokenIDs map[string][]*verifiable.VINSubject) map[int64]struct{} {
	confirmedVINs := make(map[int64]struct{})
	for vin, subjects := range vinToTokenIDs {
		if len(subjects) == 0 {
			// we don't expect this to happen
			v.logger.Warn().Str("vin", vin).Msg("no tokenIDs associated with VIN")
			continue
		}
		if len(subjects) == 1 {
			// Only one tokenID associated with this VIN
			confirmedVINs[int64(subjects[0].VehicleTokenID)] = struct{}{}
			continue
		}

		// this VIN has multiple associated tokenIDs, so only keep the one with the latest recordedAt
		v.logger.Debug().Str("vin", vin).Any("vehicleTokenIdCount", len(subjects)).Msg("VIN has multiple associated tokenIDs")
		lastTokenID := getLatestTokenID(subjects)
		confirmedVINs[lastTokenID] = struct{}{}

		// log the tokenIds of the other VCs
		if v.logger.GetLevel() <= zerolog.DebugLevel {
			for _, badSubject := range subjects {
				if badSubject.VehicleTokenID == uint32(lastTokenID) {
					continue
				}
				v.logger.Debug().Uint32("vehicleTokenId", badSubject.VehicleTokenID).
					Time("recordedAt", badSubject.RecordedAt).Str("vin", badSubject.VehicleIdentificationNumber).
					Msg("removing non latest tokenId associated with VIN")
			}
		}
	}

	return confirmedVINs
}

// getLatestTokenID returns the latest vehicleTokenId based on the recordedAt and vehicleTokenId.
func getLatestTokenID(subs []*verifiable.VINSubject) int64 {
	maxSub := subs[0]
	for _, sub := range subs {
		if maxSub.RecordedAt.Before(sub.RecordedAt) {
			maxSub = sub
		} else if maxSub.RecordedAt.Equal(sub.RecordedAt) && maxSub.VehicleTokenID < sub.VehicleTokenID {
			maxSub = sub
		}
	}
	return int64(maxSub.VehicleTokenID)
}

// isValidVC checks if a VIN credential is valid based on expiration and recording criteria.
// A credential is valid if:
// 1. It is not expired.
// 2. It was either recorded in the provided week OR recorded by a trusted source.
func (v *VINVCService) isValidVC(logger *zerolog.Logger, vehicleTokenID uint32, cred *verifiable.Credential, subject *verifiable.VINSubject, weekNum int) bool {
	if vehicleTokenID != subject.VehicleTokenID {
		logger.Warn().Uint32("vehicleTokenId", vehicleTokenID).Uint32("vcVehicleTokenId", subject.VehicleTokenID).Msg("tokenId mismatch")
		return false
	}
	// Check expiration - credential should not be expired
	startOfWeek := date.NumToWeekStart(weekNum)
	endOfWeek := date.NumToWeekEnd(weekNum)
	expiresAt, err := time.Parse(time.RFC3339, cred.ValidTo)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse ValidTo date")
		return false
	}

	if startOfWeek.After(expiresAt) {
		// Credential expired before the start of the week
		logger.Info().Str("validTo", cred.ValidTo).Msg("VINVC is expired")
		return false
	}

	// check if recorder is trusted
	isTrustedRecorder := slices.Contains(v.trustedRecorders, subject.RecordedBy)

	// credential is valid if it was recorded in the past week OR by a trusted recorder
	recordedThisWeek := false
	if !subject.RecordedAt.IsZero() {
		recordedThisWeek = subject.RecordedAt.After(startOfWeek) && subject.RecordedAt.Before(endOfWeek)
	}

	return recordedThisWeek || isTrustedRecorder
}
