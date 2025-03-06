// Package attestation provides functions for managing attestations in the context of DIMO Rewards.
package attestation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/DIMO-Network/attestation-api/pkg/grpc"
	"github.com/DIMO-Network/attestation-api/pkg/verifiable"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/services/vinvc"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Service contains all clients and services needed for attestation operations.
type Service struct {
	chClient     *ch.Client
	vinvcSrv     *vinvc.VINVCService
	attestClient pb.AttestationServiceClient
	grpcConn     *grpc.ClientConn
	settings     *config.Settings
	logger       zerolog.Logger
}

// NewService initializes a new attestation Service with all required dependencies.
func NewService(settings *config.Settings, logger *zerolog.Logger, chClient *ch.Client, vinvcSrv *vinvc.VINVCService) (*Service, error) {
	// Set up connection to server
	conn, err := grpc.NewClient(settings.AttestationAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Attestation API gRPC server: %w", err)
	}

	attestClient := pb.NewAttestationServiceClient(conn)

	return &Service{
		chClient:     chClient,
		vinvcSrv:     vinvcSrv,
		attestClient: attestClient,
		grpcConn:     conn,
		settings:     settings,
		logger:       logger.With().Str("component", "attestation-service").Logger(),
	}, nil
}

// Close closes any connections maintained by the service.
func (s *Service) Close() {
	if s.grpcConn != nil {
		s.grpcConn.Close()
	}
}

// EnsureAttestations ensures all attestations are present for the specified week.
func (s *Service) EnsureAttestations(ctx context.Context, weekNum int) error {
	weekStart := date.NumToWeekStart(weekNum)
	weekEnd := date.NumToWeekEnd(weekNum)

	activeVehicles, err := s.chClient.DescribeActiveDevices(ctx, weekStart, weekEnd)
	if err != nil {
		return fmt.Errorf("failed to get active devices: %w", err)
	}
	_, err = s.GetConfirmedVINVCs(ctx, activeVehicles, weekNum)
	if err != nil {
		return fmt.Errorf("failed to ensure VIN VCs: %w", err)
	}
	return nil
}

// GetConfirmedVINVCs performs the same actions as vinvc.Service.GetConfirmedVINVCs, but if we fail to fetch the latest valid VINVCs, we will attempt to create them.
func (s *Service) GetConfirmedVINVCs(ctx context.Context, activeVehicles []*ch.Vehicle, weekNum int) (map[int64]struct{}, error) {
	forceCreate := true
	validVinToSubjects, err := s.vinvcSrv.GetLatestValidVINVCs(ctx, activeVehicles, weekNum)
	if err != nil {
		// this is a non-fatal error, we can continue without this data
		// we will not force create VIN VCs in this case since we do not know if they are valid or not
		s.logger.Warn().Err(err).Msg("failed to get latest valid VINVCs.")
		validVinToSubjects = map[string][]*verifiable.VINSubject{}
		forceCreate = false
	}

	// convert this to a map we can use to check if a vehicle already has a VC
	validTokenIDs := make(map[int64]struct{})
	for _, subjects := range validVinToSubjects {
		for _, subject := range subjects {
			validTokenIDs[int64(subject.VehicleTokenID)] = struct{}{}
		}
	}

	// ensure VIN VCs for all active vehicles that do not have a valid VC
	for _, device := range activeVehicles {
		if _, ok := validTokenIDs[device.TokenID]; ok {
			continue
		}
		vinSubject, err := s.ensureVinVC(ctx, device.TokenID, date.NumToWeekEnd(weekNum), forceCreate)
		if err != nil {
			s.logger.Error().Err(err).Int64("vehicleTokenId", device.TokenID).Msg("Failed to ensure VIN VC")
			continue
		}
		subs := validVinToSubjects[vinSubject.VehicleIdentificationNumber]
		subs = append(subs, vinSubject)
		validVinToSubjects[vinSubject.VehicleIdentificationNumber] = subs
	}

	return s.vinvcSrv.ResolveVINConflicts(validVinToSubjects), nil
}

func (s *Service) ensureVinVC(ctx context.Context, tokenID int64, before time.Time, force bool) (*verifiable.VINSubject, error) {
	resp, err := s.attestClient.EnsureVinVc(ctx, &pb.EnsureVinVcRequest{
		TokenId: uint32(tokenID),
		Force:   force,
		// Before:  timestamppb.New(before), // TODO (kevin): need to add this to the gRPC service for when we run for a past week
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ensure VIN VC: %w", err)
	}

	var cred verifiable.Credential
	if err := json.Unmarshal([]byte(resp.GetRawVc()), &cred); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VIN VC: %w", err)
	}

	var vinSubject verifiable.VINSubject
	if err := json.Unmarshal(cred.CredentialSubject, &vinSubject); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VIN VC subject: %w", err)
	}

	return &vinSubject, nil
}
