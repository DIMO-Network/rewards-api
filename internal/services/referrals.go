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
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
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
		ContractAddress: common.HexToAddress(settings.ReferralContractAddress),
		Week:            week,
		Logger:          logger,
		UsersClient:     userClient,
	}
}

// CollectReferrals returns address pairs for referrals completed in the given week.
// These will come from referees who are earning for the first time and have a referrer
// attached to their account.
func (c *ReferralsClient) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	var res []struct {
		UserID string `boil:"r1.user_id"`
	}

	var rCols = models.RewardColumns

	err := queries.Raw(
		"SELECT DISTINCT ON (r1."+rCols.UserEthereumAddress+")"+
			" r1."+rCols.UserID+
			" FROM "+models.TableNames.Rewards+" r1"+
			" LEFT OUTER JOIN "+models.TableNames.Rewards+" r2 ON r1."+rCols.UserEthereumAddress+" = r2."+rCols.UserEthereumAddress+" AND r2."+rCols.IssuanceWeekID+" < $1"+
			" LEFT JOIN  "+models.TableNames.Vins+" v on r1."+rCols.UserDeviceTokenID+" = v."+models.VinColumns.FirstEarningTokenID+
			" WHERE r1."+rCols.IssuanceWeekID+" = $1 AND r1."+rCols.UserEthereumAddress+" IS NOT NULL AND r2."+rCols.UserDeviceID+" IS NULL"+
			" AND v."+models.VinColumns.FirstEarningWeek+" = $1"+
			" ORDER BY r1."+rCols.UserEthereumAddress+", r1."+rCols.UserID,
		issuanceWeek,
	).Bind(ctx, c.TransferService.db.DBS().Reader, &res)
	if err != nil {
		return Referrals{}, err
	}

	for _, r := range res {
		user, err := c.UsersClient.GetUser(ctx, &pb.GetUserRequest{Id: r.UserID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				c.Logger.Info().Str("userId", r.UserID).Msg("User was new this week but deleted their account.")
				continue
			}
			return refs, err
		}

		if user.ReferredBy == nil {
			continue
		}

		if user.EthereumAddress == nil {
			c.Logger.Info().Str("userId", r.UserID).Msg("Referred user does not have a valid ethereum address.")
			continue
		}

		if !user.ReferredBy.ReferrerValid {
			c.Logger.Info().Str("userId", r.UserID).Msg("Referring user has deleted their account or no longer has a confirmed ethereum address.")
			// referring eth addr is set to the referrals contract
			user.ReferredBy.EthereumAddress = c.ContractAddress[:]
		}

		refereeAddr := common.HexToAddress(*user.EthereumAddress)
		referrerAddr := common.BytesToAddress(user.ReferredBy.EthereumAddress)

		if refereeAddr == referrerAddr {
			c.Logger.Info().Str("userId", r.UserID).Msg("Referred users ethereum address is same as referring users.")
			continue
		}

		refs.Referees = append(refs.Referees, refereeAddr)
		refs.Referrers = append(refs.Referrers, referrerAddr)
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

func (c *ReferralsClient) transfer(ctx context.Context, refs Referrals) error {
	for i := 0; i < len(refs.Referees); i += c.TransferService.batchSize {
		if err := func() error {
			reqID := ksuid.New().String()
			j := i + c.TransferService.batchSize
			if j > len(refs.Referees) {
				j = len(refs.Referees)
			}

			referreesBatch := refs.Referees[i:j]
			referrersBatch := refs.Referrers[i:j]

			tx, err := c.TransferService.db.DBS().Writer.BeginTx(ctx, nil)
			if err != nil {
				return err
			}

			defer tx.Rollback() //nolint

			mtr := models.MetaTransactionRequest{
				ID:     reqID,
				Status: models.MetaTransactionRequestStatusUnsubmitted,
			}
			if err := mtr.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}

			for n := range referreesBatch {
				r := models.Referral{
					Referee:   referreesBatch[n].Bytes(),
					Referrer:  referrersBatch[n].Bytes(),
					RequestID: reqID,
				}
				if err := r.Insert(ctx, tx, boil.Infer()); err != nil {
					return err
				}
			}

			if err := c.BatchTransferReferralBonuses(reqID, referreesBatch, referrersBatch); err != nil {
				return err
			}

			if err := tx.Commit(); err != nil {
				return err
			}

			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func (c *ReferralsClient) BatchTransferReferralBonuses(requestID string, referrees []common.Address, referrers []common.Address) error {
	abi, err := contracts.ReferralMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("sendReferralBonuses", referrees, referrers)
	if err != nil {
		return err
	}
	return c.sendRequest(requestID, data)
}

func (c *ReferralsClient) sendRequest(requestID string, data []byte) error {
	event := shared.CloudEvent[transferData]{
		ID:          requestID,
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     c.ContractAddress.Hex(),
		Time:        time.Now(),
		Type:        "zone.dimo.referrals.request",
		Data: transferData{
			ID:   requestID,
			To:   c.ContractAddress,
			Data: data,
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = c.TransferService.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: c.TransferService.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	return err
}
