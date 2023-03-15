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
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReferralsClient struct {
	TransferService *TransferService
	ContractAddress common.Address
	Week            int
	Logger          *zerolog.Logger
	UsersClient     pb.UserServiceClient
}

type Referrals struct {
	Referees  []common.Address
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

// CollectReferrals returns address pairs for referrals completed in the given week.
// These will come from referees who are earning for the first time and have a referrer
// attached to their account.
func (rc *ReferralsClient) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []models.Reward

	err := models.NewQuery(
		qm.Distinct("r1."+models.RewardColumns.UserID+", r1."+models.RewardColumns.UserEthereumAddress),
		qm.From(models.TableNames.Rewards+" r1"),
		qm.LeftOuterJoin(models.TableNames.Rewards+" r2 ON r1."+models.RewardColumns.UserID+" = r2."+models.RewardColumns.UserID+" AND r2."+models.RewardColumns.IssuanceWeekID+" < ?", issuanceWeek),
		qm.Where("r1."+models.RewardColumns.IssuanceWeekID+" = ?", issuanceWeek),
		qm.Where("r2."+models.RewardColumns.UserID+" IS NULL"),
	).Bind(ctx, rc.TransferService.db.DBS().Reader, &res)
	if err != nil {
		return refs, err
	}

	for _, usr := range res {
		user, err := rc.UsersClient.GetUser(ctx, &pb.GetUserRequest{Id: usr.UserID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				rc.Logger.Info().Str("userId", usr.UserID).Msg("User was new this week but deleted their account.")
				continue
			}
			return refs, err
		}

		if user.ReferredBy == nil {
			continue
		}

		if user.EthereumAddress == nil {
			rc.Logger.Info().Str("userId", usr.UserID).Msg("Referred user does not have a valid ethereum address.")
			continue
		}

		refs.Referees = append(refs.Referees, common.HexToAddress(*user.EthereumAddress))
		refs.Referrers = append(refs.Referrers, common.BytesToAddress(user.ReferredBy.EthereumAddress))
	}

	return refs, nil
}

func (c *ReferralsClient) ReferralsIssuance(ctx context.Context) error {

	refs, err := c.CollectReferrals(ctx, c.Week)
	if err != nil {
		return err
	}

	c.Logger.Info().Msgf("Sending transactions for %d referrals.", len(refs.Referees))

	err = c.transfer(ctx, refs)
	if err != nil {
		return err
	}
	return nil
}

func (rc *ReferralsClient) transfer(ctx context.Context, refs Referrals) error {
	for i := 0; i < len(refs.Referees); i += rc.TransferService.batchSize {
		reqID := ksuid.New().String()
		j := i + rc.TransferService.batchSize
		if j > len(refs.Referees) {
			j = len(refs.Referees)
		}

		referreesBatch := refs.Referees[i:j]
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

func (rc *ReferralsClient) BatchTransferReferralBonuses(requestID string, referrees []common.Address, referrers []common.Address) error {
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
		Subject:     rc.ContractAddress.Hex(),
		Time:        time.Now(),
		Type:        "zone.dimo.referrals.request",
		Data: transferData{
			ID:   requestID,
			To:   rc.ContractAddress.Hex(),
			Data: hexutil.Encode(data),
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = rc.TransferService.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: rc.TransferService.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	return err
}
