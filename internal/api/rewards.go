package api

import (
	"context"
	"math/big"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/rewards"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

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

type averageTokens struct {
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

func (s *rewardsService) GetAverageTokens(ctx context.Context, _ *emptypb.Empty) (*pb.AverageTokensResponse, error) {
	var avrg averageTokens
	mw := make([]maxWeek, 0)
	err := models.NewQuery(qm.Select("max(issuance_week_id) as max_week"), qm.From(models.TableNames.Rewards)).Bind(ctx, s.dbs().Reader, &mw)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get max week for average tokens allocated.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	err = queries.Raw("SELECT ((sum(tokens)/ count(distinct user_device_id)) / $1::numeric)::int as average_tokens FROM rewards WHERE issuance_week_id = $2",
		ether.String(), mw[len(mw)-1].MaxWeek).Bind(ctx, s.dbs().Reader, &avrg)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get average tokens allocated for current week.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	out := &pb.AverageTokensResponse{
		AverageTokens: avrg.AverageTokens,
	}
	return out, nil
}
