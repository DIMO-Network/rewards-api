package fetchapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DIMO-Network/cloudevent"
	pb "github.com/DIMO-Network/fetch-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FetchAPIService is a service for interacting with the Fetch API.
type FetchAPIService struct {
	client pb.FetchServiceClient
	logger zerolog.Logger
}

// New creates a new instance of FetchAPIService.
func New(settings *config.Settings, logger *zerolog.Logger) (*FetchAPIService, error) {
	conn, err := grpc.NewClient(settings.FetchAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create to Fetch API gRPC client: %w", err)
	}

	return &FetchAPIService{
		client: pb.NewFetchServiceClient(conn),
		logger: logger.With().Str("component", "fetch-api-service").Logger(),
	}, nil
}

// ListCloudEvents retrieves cloud events matching the provided search criteria.
func (f *FetchAPIService) ListCloudEvents(ctx context.Context, searchOpts *pb.SearchOptions, limit int32) ([]cloudevent.CloudEvent[json.RawMessage], error) {
	resp, err := f.client.ListCloudEvents(ctx, &pb.ListCloudEventsRequest{
		Options: searchOpts,
		Limit:   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest cloud event: %w", err)
	}
	events := make([]cloudevent.CloudEvent[json.RawMessage], len(resp.GetCloudEvents()))
	for i, ce := range resp.GetCloudEvents() {
		events[i] = ce.AsCloudEvent()
	}

	return events, nil
}
