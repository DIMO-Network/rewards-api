package api

import (
	"context"

	"github.com/DIMO-Network/rewards-api/internal/database"
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
	dbs    func() *database.DBReaderWriter
	logger *zerolog.Logger
}

type totalResp struct {
	TotalPoints int64 `boil:"total_points"`
}

type totelTokens struct {
	AverageTokens int64 `boil:"average_tokens"`
}

type maxWeek struct {
	MaxWeek int64 `boil:"max_week"`
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

func (s *rewardsService) GetAverageTokens(ctx context.Context, _ *emptypb.Empty) (*pb.GetTotalPointsResponse, error) {
	tt := new(totelTokens)
	mw := new(maxWeek)
	weekQuery := models.NewQuery(
		qm.Select("max("+models.RewardColumns.IssuanceWeekID+") as max_week"),
		qm.From(models.TableNames.Rewards),
	)
	if err := weekQuery.Bind(ctx, s.dbs().Reader, mw); err != nil {
		s.logger.Err(err).Msg("Failed to get total points.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	tokenQuery := models.NewQuery(
		qm.Select("sum("+models.RewardColumns.Tokens+") / count(distinct("+models.RewardColumns.UserDeviceTokenID+")) as total_tokens"),
		qm.From(models.TableNames.Rewards),
		qm.Where(models.RewardColumns.IssuanceWeekID+" = ", mw.MaxWeek),
	)
	if err := tokenQuery.Bind(ctx, s.dbs().Reader, tt); err != nil {
		s.logger.Err(err).Msg("Failed to get total tokens.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	out := &pb.GetTotalPointsResponse{
		TotalPoints: tt.AverageTokens,
	}
	return out, nil
}
