package services

import (
	"context"

	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReferralsTask issues referral bonuses for users who
type ReferralsTask struct {
	Logger      *zerolog.Logger
	UsersClient pb.UserServiceClient
	DB          db.Store
}

type Referrals struct {
	Referees []common.Address
	Referrer []common.Address
}

// CollectReferrals returns address information for referrals completed in the given week.
// These will come from referees who are earning for the first time and have a referrer
// attached to their account.
func (r *ReferralsTask) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []models.Reward

	err := models.NewQuery(
		qm.Distinct("r1."+models.RewardColumns.UserID+", r1."+models.RewardColumns.UserEthereumAddress),
		qm.From(models.TableNames.Rewards+" r1"),
		qm.LeftOuterJoin(models.TableNames.Rewards+" r2 ON r1."+models.RewardColumns.UserID+" = r2."+models.RewardColumns.UserID+" AND r2."+models.RewardColumns.IssuanceWeekID+" < ?", issuanceWeek),
		qm.Where("r1."+models.RewardColumns.IssuanceWeekID+" = ?", issuanceWeek),
		qm.Where("r2."+models.RewardColumns.UserID+" IS NULL"),
	).Bind(ctx, r.DB.DBS().Reader, &res)

	// err := queries.Raw(`SELECT DISTINCT r1.user_id, r1.user_ethereum_address FROM rewards r1 LEFT
	// OUTER JOIN rewards r2 ON r1.user_id = r2.user_id AND r2.issuance_week_id < $1 WHERE
	// r1.issuance_week_id = $1 AND r2.user_id IS NULL`, issuanceWeek)
	if err != nil {
		return refs, err
	}

	for _, usr := range res {

		user, err := r.UsersClient.GetUser(ctx, &pb.GetUserRequest{
			Id: usr.UserID,
		})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				r.Logger.Info().Msg("User has deleted their account.")
				continue
			}
			return refs, err
		}

		if user.ReferredBy == nil {
			continue
		}

		refs.Referees = append(refs.Referees, common.HexToAddress(*user.EthereumAddress))
		refs.Referrer = append(refs.Referrer, common.BytesToAddress(user.ReferredBy.EthereumAddress))
	}

	return refs, nil
}
