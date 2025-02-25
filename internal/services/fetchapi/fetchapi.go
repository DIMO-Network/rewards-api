package fetchapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	pb "github.com/DIMO-Network/fetch-api/pkg/grpc"
	"github.com/DIMO-Network/model-garage/pkg/cloudevent"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FetchAPIService is a service for interacting with the Fetch API.
type FetchAPIService struct {
	fetchGRPCAddr string
	vehicleAddr   common.Address
	chainID       uint64
	client        pb.FetchServiceClient
	once          sync.Once
	logger        zerolog.Logger
}

// New creates a new instance of FetchAPIService
func New(settings *config.Settings, logger *zerolog.Logger) *FetchAPIService {
	return &FetchAPIService{
		fetchGRPCAddr: settings.FetchAPIGRPCEndpoint,
		vehicleAddr:   settings.VehicleNFTAddress,
		chainID:       uint64(settings.ChainID),
		logger:        logger.With().Str("component", "fetch_api_service").Logger(),
	}
}

// GetLatestCloudEvent retrieves the most recent cloud event matching the provided search criteria
func (f *FetchAPIService) GetLatestCloudEvent(ctx context.Context, filter *pb.SearchOptions) (cloudevent.CloudEvent[json.RawMessage], error) {
	client, err := f.getClient()
	if err != nil {
		return cloudevent.CloudEvent[json.RawMessage]{}, fmt.Errorf("failed to initialize gRPC client: %w", err)
	}

	resp, err := client.GetLatestCloudEvent(ctx, &pb.GetLatestCloudEventRequest{
		Options: filter,
	})
	if err != nil {
		return cloudevent.CloudEvent[json.RawMessage]{}, fmt.Errorf("failed to get latest cloud event: %w", err)
	}

	return resp.GetCloudEvent().AsCloudEvent(), nil
}

// getClient initializes the gRPC client if not already initialized.
func (f *FetchAPIService) getClient() (pb.FetchServiceClient, error) {
	if f.client != nil {
		return f.client, nil
	}
	var err error
	f.once.Do(func() {
		var conn *grpc.ClientConn
		conn, err = grpc.NewClient(f.fetchGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			err = fmt.Errorf("failed to connect to Fetch API gRPC server: %w", err)
			return
		}
		f.client = pb.NewFetchServiceClient(conn)
	})
	return f.client, err
}

// CreateVehicleDID creates a DID for a vehicle.
func (f *FetchAPIService) CreateVehicleDID(tokenID uint32) string {
	return cloudevent.NFTDID{
		ChainID:         f.chainID,
		ContractAddress: f.vehicleAddr,
		TokenID:         tokenID,
	}.String()
}
