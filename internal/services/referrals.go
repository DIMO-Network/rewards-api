package services

import (
	"context"

	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// ReferralsTask controller to collect referrals and issue bonus
type ReferralsTask struct {
	Logger          *zerolog.Logger
	UsersClient     pb_users.UserServiceClient
	DataService     DeviceActivityClient
	DB              db.Store
	TransferService Transfer
}

func (r *ReferralsTask) CollectReferrals(issuanceWeek int) ([]string, error) {
	ctx := context.Background()

	historicalUsers, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.NEQ(issuanceWeek),
		qm.Select(models.RewardColumns.UserID)).All(ctx, r.DB.DBS().Reader)
	if err != nil {
		return []string{}, err
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
		return []string{}, err
	}

	newUsers := make([]string, len(newUserResp))
	for n, usr := range newUserResp {
		newUsers[n] = usr.UserID
	}

	return newUsers, nil
}
