package services

import (
	"context"

	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

// ReferralsTask controller to collect referrals and issue bonus
type ReferralsTask struct {
	Logger          *zerolog.Logger
	UsersClient     pb_users.UserServiceClient
	DataService     DeviceActivityClient
	DB              db.Store
	TransferService Transfer
}

type Referrals struct {
	Referred common.Address
	Referrer common.Address
}

// CollectReferrals Check if users who recieved rewards for the first time this week were referred
// if they were, collect their address and the address of their referrer
func (r *ReferralsTask) CollectReferrals(issuanceWeek int) ([]Referrals, error) {
	ctx := context.Background()

	historicalUsers, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.NEQ(issuanceWeek)).All(ctx, r.DB.DBS().Reader)
	if err != nil {
		return []Referrals{}, err
	}

	historicalUserIDs := make([]string, len(historicalUsers))
	// historicalEthAddrs := make([]string, len(historicalUsers))

	for _, usr := range historicalUsers {
		historicalUserIDs = append(historicalUserIDs, usr.UserID)
		// historicalEthAddrs = append(historicalEthAddrs, usr.UserEthereumAddress.String)
	}

	newUserResp, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
		models.RewardWhere.UserID.NIN(historicalUserIDs),
		// models.RewardWhere.UserEthereumAddress.NIN(historicalEthAddrs),
	).
		All(ctx, r.DB.DBS().Reader)
	if err != nil {
		return []Referrals{}, err
	}

	refs := make([]Referrals, 0)
	for _, usr := range newUserResp {

		user, err := r.UsersClient.GetUser(ctx, &pb_users.GetUserRequest{
			Id: usr.UserID,
		})
		if err != nil {
			return []Referrals{}, err
		}

		if user.ReferredBy == nil {
			continue
		}

		refs = append(refs, Referrals{
			Referred: common.HexToAddress(*user.EthereumAddress),
			Referrer: common.BytesToAddress(user.ReferredBy.EthereumAddress),
		})
	}

	return refs, nil
}
