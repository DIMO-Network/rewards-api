package services

import (
	"context"

	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReferralsTask controller to collect referrals and issue bonus
type ReferralsTask struct {
	Logger          *zerolog.Logger
	UsersClient     pb.UserServiceClient
	DataService     DeviceActivityClient
	DB              db.Store
	TransferService Transfer
}

type Referrals struct {
	Referred []common.Address
	Referrer []common.Address
}

// CollectReferrals Check if users who recieved rewards for the first time this week were referred
// if they were, collect their address and the address of their referrer
func (r *ReferralsTask) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []models.Reward

	err := queries.Raw(`SELECT DISTINCT r1.user_id, r1.user_ethereum_address FROM rewards r1 LEFT
	OUTER JOIN rewards r2 ON r1.user_id = r2.user_id AND r2.issuance_week_id < $1 WHERE
	r1.issuance_week_id = $1 AND r2.user_id IS NULL`, issuanceWeek).Bind(ctx, r.DB.DBS().Reader, &res)
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

		refs.Referred = append(refs.Referred, common.HexToAddress(*user.EthereumAddress))
		refs.Referrer = append(refs.Referrer, common.BytesToAddress(user.ReferredBy.EthereumAddress))
	}

	return refs, nil
}
