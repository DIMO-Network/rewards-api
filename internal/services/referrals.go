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
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReferralsTask to collect referrals and issue bonuses
type ReferralsTask struct {
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
	Referees  []common.Address
	Referrers []common.Address
}

func NewRewardBonusTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	contractAddress common.Address,
	usersClient pb.UserServiceClient,
	db db.Store,
	logger *zerolog.Logger) *ReferralsTask {

	return &ReferralsTask{
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

// CollectReferrals returns address pairs for referrals completed in the given week.
// These will come from referees who are earning for the first time and have a referrer
// attached to their account.
func (t *ReferralsTask) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []models.Reward

	err := models.NewQuery(
		qm.Distinct("r1."+models.RewardColumns.UserID+", r1."+models.RewardColumns.UserEthereumAddress),
		qm.From(models.TableNames.Rewards+" r1"),
		qm.LeftOuterJoin(models.TableNames.Rewards+" r2 ON r1."+models.RewardColumns.UserID+" = r2."+models.RewardColumns.UserID+" AND r2."+models.RewardColumns.IssuanceWeekID+" < ?", issuanceWeek),
		qm.Where("r1."+models.RewardColumns.IssuanceWeekID+" = ?", issuanceWeek),
		qm.Where("r2."+models.RewardColumns.UserID+" IS NULL"),
	).Bind(ctx, t.db.DBS().Reader, &res)
	if err != nil {
		return refs, err
	}

	for _, usr := range res {
		user, err := t.UsersClient.GetUser(ctx, &pb.GetUserRequest{Id: usr.UserID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				t.Logger.Info().Str("userId", usr.UserID).Msg("User was new this week but deleted their account.")
				continue
			}
			return refs, err
		}

		if user.ReferredBy == nil {
			continue
		}

		// TODO(elffjs): What if this is nil?
		refs.Referees = append(refs.Referees, common.HexToAddress(*user.EthereumAddress))
		refs.Referrers = append(refs.Referrers, common.BytesToAddress(user.ReferredBy.EthereumAddress))
	}

	return refs, nil
}

func (t *ReferralsTask) TransferReferralBonuses(ctx context.Context, refs Referrals) error {
	for i := 0; i < len(refs.Referees); i += t.batchSize {
		reqID := ksuid.New().String()
		j := i + t.batchSize
		if j > len(refs.Referees) {
			j = len(refs.Referees)
		}
		err := t.BatchTransferReferralBonuses(reqID, refs.Referees[i:j], refs.Referrers[i:j])
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *ReferralsTask) BatchTransferReferralBonuses(requestID string, referreds []common.Address, referrers []common.Address) error {
	abi, err := contracts.ReferralMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("sendReferralBonuses", referreds, referrers)
	if err != nil {
		return err
	}
	return t.sendRequest(requestID, data)
}

func (t *ReferralsTask) sendRequest(requestID string, data []byte) error {
	event := shared.CloudEvent[transferData]{
		ID:          requestID,
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     hexutil.Encode(t.ContractAddress[:]),
		Time:        time.Now(),
		Type:        "zone.dimo.transaction.request",
		Data: transferData{
			ID:   requestID,
			To:   t.ContractAddress,
			Data: data,
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = t.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: t.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	return err
}
