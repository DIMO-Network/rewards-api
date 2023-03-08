package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/sqlboiler/queries"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReferralsClient controller to collect referrals and issue bonus
type ReferralsClient struct {
	Producer        sarama.SyncProducer
	Consumer        sarama.ConsumerGroup
	RequestTopic    string
	StatusTopic     string
	db              db.Store
	ContractAddress common.Address
	batchSize       int
	UsersClient     pb.UserServiceClient
	Logger          *zerolog.Logger
}

type Referrals struct {
	Referreds []common.Address
	Referrers []common.Address
}

func NewRewardBonusTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	contractAddress common.Address,
	usersClient pb.UserServiceClient,
	db db.Store,
	logger *zerolog.Logger) *ReferralsClient {

	return &ReferralsClient{
		Producer:        producer,
		RequestTopic:    settings.MetaTransactionSendTopic,
		StatusTopic:     settings.MetaTransactionStatusTopic,
		db:              db,
		ContractAddress: contractAddress,
		batchSize:       settings.TransferBatchSize,
		UsersClient:     usersClient,
		Logger:          logger,
	}
}

// CollectReferrals Check if users who recieved rewards for the first time this week were referred
// if they were, collect their address and the address of their referrer
func (rc *ReferralsClient) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []models.Reward

	err := queries.Raw(`SELECT DISTINCT r1.user_id, r1.user_ethereum_address FROM rewards r1 LEFT
	OUTER JOIN rewards r2 ON r1.user_id = r2.user_id AND r2.issuance_week_id < $1 WHERE
	r1.issuance_week_id = $1 AND r2.user_id IS NULL`, issuanceWeek).Bind(ctx, rc.db.DBS().Reader, &res)
	if err != nil {
		return refs, err
	}

	for _, usr := range res {

		user, err := rc.UsersClient.GetUser(ctx, &pb.GetUserRequest{
			Id: usr.UserID,
		})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				rc.Logger.Info().Msg("User has deleted their account.")
				continue
			}
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

func (rc *ReferralsClient) TransferReferralBonuses(ctx context.Context, refs Referrals) error {
	err := rc.transferReferralBonuses(ctx, refs)
	if err != nil {
		return err
	}
	return nil
}

func (rc *ReferralsClient) transferReferralBonuses(ctx context.Context, refs Referrals) error {

	for i := 0; i < len(refs.Referreds); i += rc.batchSize {
		reqID := ksuid.New().String()
		j := i + rc.batchSize
		if j > len(refs.Referreds) {
			j = len(refs.Referreds)
		}
		err := rc.BatchTransferReferralBonuses(reqID, refs.Referreds[i:j], refs.Referrers[i:j])
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *ReferralsClient) BatchTransferReferralBonuses(requestID string, referreds []common.Address, referrers []common.Address) error {
	abi, err := contracts.ReferralsMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("sendReferralBonuses", referreds, referrers)
	if err != nil {
		return err
	}
	return rc.sendRequest(requestID, data)
}

func (rc *ReferralsClient) sendRequest(requestID string, data []byte) error {
	event := shared.CloudEvent[transferData]{
		ID:          requestID,
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     hexutil.Encode(rc.ContractAddress[:]),
		Time:        time.Now(),
		Type:        "zone.dimo.referrals.request",
		Data: transferData{
			ID:   requestID,
			To:   hexutil.Encode(rc.ContractAddress[:]),
			Data: hexutil.Encode(data),
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = rc.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: rc.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	return err
}
