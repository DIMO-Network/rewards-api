package api

import (
	"context"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/rewards"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewRewardsService(dbs func() *database.DBReaderWriter, logger *zerolog.Logger) pb.RewardsServiceServer {
	return &rewardsService{
		dbs:    dbs,
		logger: logger,
	}
}

type rewardsService struct {
	pb.UnimplementedRewardsServiceServer
	dbs        func() *database.DBReaderWriter
	logger     *zerolog.Logger
	dataClient services.DeviceDataClient
}

type totalResp struct {
	TotalPoints int64 `boil:"total_points"`
}

func (s *rewardsService) GetTotalPoints(ctx context.Context, _ *emptypb.Empty) (*pb.GetTotalPointsResponse, error) {
	tp := new(totalResp)
	query := models.NewQuery(
		qm.Select("sum("+models.RewardColumns.StreakPoints+" + "+models.RewardColumns.IntegrationPoints+") as total_points"),
		qm.From(models.TableNames.Rewards),
	)
	if err := query.Bind(ctx, s.dbs().Reader, tp); err != nil {
		s.logger.Err(err).Msg("Failed to get total points.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	out := &pb.GetTotalPointsResponse{
		TotalPoints: tp.TotalPoints,
	}
	return out, nil
}

func (s *rewardsService) GetQualifiedDevices(ctx context.Context, req *pb.GetQualifiedDevicesRequest) (*pb.GetQualifiedDevicesResponse, error) {
	start, end := req.Start.AsTime(), req.End.AsTime()
	data, err := s.dataClient.DescribeActiveDevices(start, end)
	if err != nil {
		return nil, err
	}

	devices := make([]*pb.GetQualifiedDevicesDevice, len(data))

	for i, dev := range data {
		devices[i] = &pb.GetQualifiedDevicesDevice{Id: dev.ID, IntegrationIds: dev.Integrations}
	}

	return &pb.GetQualifiedDevicesResponse{Devices: devices}, nil
}
