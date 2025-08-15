package api

import (
	"context"
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	pb "github.com/DIMO-Network/shared/api/rewards"
	"github.com/DIMO-Network/shared/pkg/db"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

func NewRewardsService(dbs db.Store, logger *zerolog.Logger) pb.RewardsServiceServer {
	return &rewardsService{
		dbs:    dbs,
		logger: logger,
	}
}

type rewardsService struct {
	pb.UnimplementedRewardsServiceServer
	dbs    db.Store
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
	if err := query.Bind(ctx, s.dbs.DBS().Reader, tp); err != nil {
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
	err := models.NewQuery(qm.Select("max(issuance_week_id) as max_week"), qm.From(models.TableNames.Rewards)).Bind(ctx, s.dbs.DBS().Reader, &mw)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get max week for average tokens allocated.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	// sum tokens allocated in the current week and divide by devices; then divide by ethers value to convert from gwei
	err = queries.Raw("SELECT ((sum(tokens)/ count(distinct user_device_id)) / $1::numeric)::int as average_tokens FROM rewards WHERE issuance_week_id = $2",
		ether.String(), mw[len(mw)-1].MaxWeek).Bind(ctx, s.dbs.DBS().Reader, &avrg)
	if err != nil {
		s.logger.Err(err).Msg("Failed to get average tokens allocated for current week.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	out := &pb.AverageTokensResponse{
		AverageTokens: avrg.AverageTokens,
	}
	return out, nil
}

func (s *rewardsService) GetDeviceRewards(ctx context.Context, req *pb.GetDeviceRewardsRequest) (*pb.GetDeviceRewardsResponse, error) {
	rs, err := models.Rewards(
		models.RewardWhere.UserDeviceID.EQ(req.Id),
		qm.OrderBy(models.RewardColumns.IssuanceWeekID+" ASC"),
	).All(ctx, s.dbs.DBS().Reader)
	if err != nil {
		s.logger.Err(err).Str("userDeviceId", req.Id).Msg("Failed to get rewards for device.")
		return nil, status.Error(codes.Internal, "Internal error.")
	}

	resp := pb.GetDeviceRewardsResponse{
		Id:     req.Id,
		Tokens: 0,
		Weeks:  []*pb.DeviceRewardsWeek{},
	}

	for _, r := range rs {
		var tokEth float64
		if !r.Tokens.IsZero() {
			tokEth, _ = r.Tokens.Float64()
			tokEth /= 1e18
		}

		resp.Tokens += tokEth

		row := pb.DeviceRewardsWeek{
			EndDate:             date.NumToWeekEnd(r.IssuanceWeekID).UTC().Format("2006-01-02"),
			Tokens:              tokEth,
			ConnectionStreak:    int32(r.ConnectionStreak),
			DisconnectionStreak: int32(r.DisconnectionStreak),
			IntegrationIds:      r.IntegrationIds,
		}
		resp.Weeks = append(resp.Weeks, &row)
	}

	return &resp, nil
}

func (s *rewardsService) GetBlacklistStatus(ctx context.Context, req *pb.GetBlacklistStatusRequest) (*pb.GetBlacklistStatusResponse, error) {
	if len(req.EthereumAddress) != common.AddressLength {
		return nil, status.Errorf(codes.InvalidArgument, "Ethereum address had length %d instead of the required %d.", len(req.EthereumAddress), common.AddressLength)
	}

	addr := common.BytesToAddress(req.EthereumAddress).Hex()

	bl, err := models.FindBlacklist(ctx, s.dbs.DBS().Reader, addr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var t timestamppb.Timestamp
			return &pb.GetBlacklistStatusResponse{
				IsBlacklisted: false,
				Note:          "",
				CreatedAt:     &t,
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "Database lookup failed: %v", err)
	}

	return &pb.GetBlacklistStatusResponse{
		IsBlacklisted: true,
		Note:          bl.Note,
		CreatedAt:     timestamppb.New(bl.CreatedAt),
	}, nil
}

var timeNow = time.Now

func (s *rewardsService) SetBlacklistStatus(ctx context.Context, req *pb.SetBlacklistStatusRequest) (*pb.SetBlacklistStatusResponse, error) {
	if len(req.EthereumAddress) != common.AddressLength {
		return nil, status.Errorf(codes.InvalidArgument, "Ethereum address had length %d instead of the required %d.", len(req.EthereumAddress), common.AddressLength)
	}

	addr := common.BytesToAddress(req.EthereumAddress).Hex()

	bl := models.Blacklist{
		UserEthereumAddress: addr, // Only field used for removal.
		CreatedAt:           timeNow(),
		Note:                req.Note,
	}

	if req.IsBlacklisted {
		// Only upserting here so that we can do nothing if the address is already in the list.
		// TODO(elffjs): Tell the caller about this case.
		if err := bl.Upsert(ctx, s.dbs.DBS().Writer, false, []string{models.BlacklistColumns.UserEthereumAddress}, boil.Infer(), boil.Infer()); err != nil {
			return nil, status.Errorf(codes.Internal, "Database insert failed: %v", err)
		}
	} else {
		// TODO(elfjjs): Tell the caller whether the address was in there to begin with.
		if _, err := bl.Delete(ctx, s.dbs.DBS().Writer); err != nil {
			return nil, status.Errorf(codes.Internal, "Deletion failed: %v", err)
		}
	}

	return &pb.SetBlacklistStatusResponse{}, nil
}
