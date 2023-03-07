package services

import (
	"context"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/queries"
)

type ReferralsClient struct {
	TransferService *TransferService
	ContractAddress common.Address
	Week            int
	Logger          *zerolog.Logger
	UsersClient     pb.UserServiceClient
}

type Referrals struct {
	Referreds []common.Address
	Referrers []common.Address
}

func NewRewardBonusTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	contractAddress common.Address,
	// settings config.Settings, producer sarama.SyncProducer, reqTopic string, contract Contract,
	transferService *TransferService,
	db db.Store) *ReferralsClient {

	return &ReferralsClient{
		ContractAddress: contractAddress,
		TransferService: transferService,
	}
}

// CollectReferrals Check if users who recieved rewards for the first time this week were referred
// if they were, collect their address and the address of their referrer
func (r *ReferralsClient) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []models.Reward

	err := queries.Raw(`SELECT DISTINCT r1.user_id, r1.user_ethereum_address FROM rewards r1 LEFT
	OUTER JOIN rewards r2 ON r1.user_id = r2.user_id AND r2.issuance_week_id < $1 WHERE
	r1.issuance_week_id = $1 AND r2.user_id IS NULL`, issuanceWeek).Bind(ctx, r.TransferService.db.DBS().Reader, &res)
	if err != nil {
		return refs, err
	}

	for _, usr := range res {

		user, err := r.UsersClient.GetUser(ctx, &pb.GetUserRequest{
			Id: usr.UserID,
		})
		if err != nil {
			return refs, err
		}

		if user.ReferredBy == nil {
			continue
		}

		refs.Referreds = append(refs.Referreds, common.HexToAddress(*user.EthereumAddress))
		refs.Referrers = append(refs.Referrers, common.BytesToAddress(user.ReferredBy.EthereumAddress))
	}

	return refs, nil
}
