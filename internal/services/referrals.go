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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ReferralsClient struct {
	TransferService *TransferService
	ContractAddress common.Address
	Week            int
	Logger          *zerolog.Logger
	UsersClient     pb.UserServiceClient
}

type Referrals struct {
	Referrees []common.Address
	Referrers []common.Address
}

func NewReferralBonusService(
	settings *config.Settings,
	transferService *TransferService,
	week int,
	logger *zerolog.Logger,
	userClient pb.UserServiceClient) *ReferralsClient {

	return &ReferralsClient{
		TransferService: transferService,
		ContractAddress: common.HexToAddress(settings.IssuanceContractAddress),
		Week:            week,
		Logger:          logger,
		UsersClient:     userClient,
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

		refs.Referrees = append(refs.Referrees, common.HexToAddress(*user.EthereumAddress))
		refs.Referrers = append(refs.Referrers, common.BytesToAddress(user.ReferredBy.EthereumAddress))
	}

	return refs, nil
}

func (c *ReferralsClient) ReferralsIssuance(ctx context.Context) error {

	refs, err := c.CollectReferrals(ctx, c.Week)
	if err != nil {
		return err
	}

	err = c.transfer(ctx, refs)
	if err != nil {
		return err
	}
	return nil
}

func (rc *ReferralsClient) transfer(ctx context.Context, refs Referrals) error {
	for i := 0; i < len(refs.Referrees); i += rc.TransferService.batchSize {
		reqID := ksuid.New().String()
		j := i + rc.TransferService.batchSize
		if j > len(refs.Referrees) {
			j = len(refs.Referrees)
		}

		referreesBatch := refs.Referrees[i:j]
		referrersBatch := refs.Referrers[i:j]
		tx, err := rc.TransferService.db.DBS().Writer.BeginTx(context.Background(), nil)
		if err != nil {
			return err
		}

		defer tx.Rollback() //nolint
		for n, user := range referreesBatch {
			r := models.Referral{
				JobStatus: models.ReferralsJobStatusStarted,
				Referred:  user[:],
				Referrer:  referrersBatch[n][:],
			}
			err := r.Insert(ctx, rc.TransferService.db.DBS().Writer, boil.Infer())
			if err != nil {
				return err
			}
		}
		err = rc.BatchTransferReferralBonuses(reqID, referreesBatch, referrersBatch)
		if err != nil {
			return err
		}
	}
	return nil
}


func (c *ReferralsClient) BatchTransferReferralBonuses(requestID string, referrees []common.Address, referrers []common.Address) error {
	abi, err := contracts.ReferralsMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("sendReferralBonuses", referrees, referrers)
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
