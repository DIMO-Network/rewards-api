package vinvc

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/DIMO-Network/attestation-api/pkg/verifiable"
	pb "github.com/DIMO-Network/fetch-api/pkg/grpc"
	"github.com/DIMO-Network/model-garage/pkg/cloudevent"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var dincSource = common.HexToAddress("0x4F098Ea7cAd393365b4d251Dd109e791e6190239")

// FetchAPIService defines the interface Fetch API.
type FetchAPIService interface {
	// GetLatestCloudEvent retrieves the most recent cloud event matching the provided search criteria
	GetLatestCloudEvent(ctx context.Context, filter *pb.SearchOptions) (cloudevent.CloudEvent[json.RawMessage], error)
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
		logger:           logger.With().Str("component", "vinvc_service").Logger(),
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
// 3. The TokenID is the only one associated with that VIN.
func (v *VINVCService) GetConfirmedVINVCs(ctx context.Context, activeVehicles []*ch.Vehicle) (map[int64]struct{}, error) {
	confirmedVINs := make(map[int64]struct{})
	// Map to track VINs and their associated tokenIDs (only for valid VCs)
	validVinToTokenIDs := make(map[string][]int64)

	// First pass: collect all valid VINVCs and their associated VINs
	for _, vehicle := range activeVehicles {
		logger := v.logger.With().Int64("vehicleTokenId", vehicle.TokenID).Logger()
		// Set search options
		opts := &pb.SearchOptions{
			DataVersion: &wrapperspb.StringValue{Value: v.vinVCDataVersion},
			Type:        &wrapperspb.StringValue{Value: cloudevent.TypeVerifableCredential},
			Subject:     &wrapperspb.StringValue{Value: cloudevent.NFTDID{ChainID: v.chainID, TokenID: uint32(vehicle.TokenID)}.String()},
		}

		// Get latest cloud event
		cloudEvent, err := v.fetchService.GetLatestCloudEvent(ctx, opts)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get latest VIN VC")
			continue
		}

		// Parse and validate the credential
		cred := verifiable.Credential{}
		if err := json.Unmarshal(cloudEvent.Data, &cred); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal VIN credential")
			continue
		}

		// Parse the credential subject
		credSubject := verifiable.VINSubject{}
		if err := json.Unmarshal(cred.CredentialSubject, &credSubject); err != nil {
			logger.Error().Err(err).Msg("failed to unmarshal VIN credential subject")
			continue
		}

		// Skip if VIN is empty
		if credSubject.VehicleIdentificationNumber == "" {
			logger.Warn().Msg("VINVC has empty VIN")
			continue
		}

		// Check if this is a valid VC
		if v.isValidVC(&cred, &credSubject) {
			// Only track VINs from valid VCs
			vin := credSubject.VehicleIdentificationNumber
			validVinToTokenIDs[vin] = append(validVinToTokenIDs[vin], vehicle.TokenID)
			confirmedVINs[vehicle.TokenID] = struct{}{}
		} else {
			v.logger.Info().Any("credentialSubject", credSubject).Msg("VINVC did not meet validation criteria")
		}
	}

	// Second pass: check for VIN uniqueness and add to confirmed list
	for vin, tokenIDs := range validVinToTokenIDs {
		if len(tokenIDs) > 1 {
			// This VIN has multiple associated tokenIDs, so none are confirmed
			v.logger.Warn().Str("vin", vin).Any("vehicleTokenIds", tokenIDs).Msg("VIN has multiple associated tokenIDs")
			for _, tokenID := range tokenIDs {
				delete(confirmedVINs, tokenID)
			}
		}
	}

	return confirmedVINs, nil
}

// isValidVC checks if a VIN credential is valid based on expiration and recording criteria.
// A credential is valid if:
// 1. It is not expired.
// 2. It was either recorded in the past week OR recorded by a trusted source.
func (v *VINVCService) isValidVC(cred *verifiable.Credential, subject *verifiable.VINSubject) bool {
	// Check expiration - credential should not be expired
	// previous Monday
	endOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	expiresAt, err := time.Parse(time.RFC3339, cred.ValidTo)
	if err != nil {
		v.logger.Error().Err(err).Msg("failed to parse ValidTo date")
		return false
	}

	if endOfWeek.After(expiresAt) {
		// Credential is expired
		v.logger.Info().Str("validTo", cred.ValidTo).Msg("VINVC is expired")
		return false
	}

	oneWeekAgo := endOfWeek.AddDate(0, 0, -7)

	// Check if recorder is trusted
	isTrustedRecorder := slices.Contains(v.trustedRecorders, subject.RecordedBy)

	// Credential is valid if it was recorded in the past week OR by a trusted recorder
	isRecentlyRecorded := false
	if !subject.RecordedAt.IsZero() {
		isRecentlyRecorded = subject.RecordedAt.After(oneWeekAgo)
	}

	return isRecentlyRecorded || isTrustedRecorder
}
