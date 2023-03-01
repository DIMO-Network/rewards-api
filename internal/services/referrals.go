package services

import (
	"context"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
)

// ReferralsTask controller to collect referrals and issue bonus
type ReferralsTask struct {
	Logger               *zerolog.Logger
	UsersClient          pb_users.UserServiceClient
	DB                   db.Store
	ReferralBonusService ReferralBonusService
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
	db db.Store) ReferralBonusService {

	return &Client{
		ContractAddress: contractAddress,
		Producer:        producer,
		RequestTopic:    settings.MetaTransactionSendTopic,
		StatusTopic:     settings.MetaTransactionStatusTopic,
		db:              db,
		batchSize:       settings.TransferBatchSize,
	}
}

type ReferralBonusService interface {
	// CollectReferrals(ctx context.Context, week int) error
	TransferReferralBonuses(ctx context.Context, refs Referrals) error
}

// CollectReferrals Check if users who recieved rewards for the first time this week were referred
// if they were, collect their address and the address of their referrer
func (r *ReferralsTask) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	historicalUsers, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.NEQ(issuanceWeek)).All(ctx, r.DB.DBS().Reader)
	if err != nil {
		return refs, err
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
		return refs, err
	}

	for _, usr := range newUserResp {

		user, err := r.UsersClient.GetUser(ctx, &pb_users.GetUserRequest{
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

func (c *Client) TransferReferralBonuses(ctx context.Context, refs Referrals) error {
	err := c.transferReferralBonuses(ctx, refs)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) transferReferralBonuses(ctx context.Context, refs Referrals) error {

	for i := 0; i < len(refs.Referreds); i += c.batchSize {
		reqID := ksuid.New().String()
		j := i + c.batchSize
		if j > len(refs.Referreds) {
			j = len(refs.Referreds)
		}
		err := c.BatchTransferReferralBonuses(reqID, refs.Referreds[i:j], refs.Referrers[i:j])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) BatchTransferReferralBonuses(requestID string, referreds []common.Address, referrers []common.Address) error {
	abi, err := contracts.ReferralMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("sendReferralBonuses", referreds, referrers)
	if err != nil {
		return err
	}
	return c.sendRequest(requestID, data)
}
