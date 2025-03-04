package attestation

import (
	"context"
	"fmt"

	pb "github.com/DIMO-Network/attestation-api/pkg/grpc"
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
	logger       *zerolog.Logger
}

// NewService initializes a new attestation Service with all required dependencies.
func NewService(settings *config.Settings, logger *zerolog.Logger, chClient *ch.Client, vinvcSrv *vinvc.VINVCService) (*Service, error) {
	// Set up connection to server
	conn, err := grpc.NewClient(settings.AttestationAPIGRPCEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		logger:       logger,
	}, nil
}

// Close closes any connections maintained by the service.
func (s *Service) Close() {
	if s.grpcConn != nil {
		s.grpcConn.Close()
	}
}

// EnsureVINVCs pulls attestation data for the specified week.
func (s *Service) EnsureVINVCs(weekNum int) error {
	weekStart := date.NumToWeekStart(weekNum)
	weekEnd := date.NumToWeekEnd(weekNum)
	ctx := context.Background()

	activeDevices, err := s.chClient.DescribeActiveDevices(ctx, weekStart, weekEnd)
	if err != nil {
		return fmt.Errorf("failed to get active devices: %w", err)
	}

	vinVCConfirmed, err := s.vinvcSrv.GetConfirmedVINVCs(ctx, activeDevices, weekNum)
	forceCreate := true
	if err != nil {
		// this is a non-fatal error, we can continue without this data
		// we will not force create VIN VCs in this case since we do not know if they are valid or not
		s.logger.Warn().Err(err).Msg("Failed to get confirmed VIN VC VINs. continuing execution.")
		vinVCConfirmed = map[int64]struct{}{}
		forceCreate = false
	}
	for _, device := range activeDevices {
		if _, ok := vinVCConfirmed[device.TokenID]; ok {
			continue
		}
		_, err := s.attestClient.EnsureVinVc(context.Background(), &pb.EnsureVinVcRequest{
			TokenId: uint32(device.TokenID),
			Force:   forceCreate,
		})
		if err != nil {
			s.logger.Error().Err(err).Int64("vehicleTokenId", device.TokenID).Msg("Failed to ensure VIN VC")
			continue
		}
	}
	return nil
}
