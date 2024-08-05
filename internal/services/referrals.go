package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type ReferralsClient struct {
	TransferService *TransferService
	ContractAddress common.Address
	Week            int
	Logger          *zerolog.Logger
	UsersClient     pb.UserServiceClient
}

type Referrals struct {
	Referees        []common.Address
	Referrers       []common.Address
	RefereeUserIDs  []string
	ReferrerUserIDs []string
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

const level2Weeks = 4

type UserTuple struct {
	ID      string
	Address common.Address
}

// CollectReferrals returns address pairs for referrals completed in the given week.
// These will come from referees who are earning for the first time and have a referrer
// attached to their account.
func (c *ReferralsClient) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	logger := c.Logger.With().Int("issuanceWeek", issuanceWeek).Logger()

	vehicleNFTsHittingLevel2, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
		models.RewardWhere.ConnectionStreak.EQ(level2Weeks),
		models.RewardWhere.DisconnectionStreak.EQ(0),
	).All(ctx, c.TransferService.db.DBS().Reader)
	if err != nil {
		return refs, nil
	}

	logger.Info().Msgf("Had %d vehicle NFTs hit level 2.", len(vehicleNFTsHittingLevel2))

	numVehiclesLevel2FirstTime := 0
	ownersOfLevel2FirstTimeVehicles := make(map[common.Address]struct{})

	for _, r := range vehicleNFTsHittingLevel2 {
		if vehicleHitLevel2Before, err := models.Rewards(
			models.RewardWhere.IssuanceWeekID.LT(issuanceWeek),
			models.RewardWhere.ConnectionStreak.EQ(level2Weeks),
			models.RewardWhere.UserDeviceTokenID.EQ(r.UserDeviceTokenID),
		).Exists(ctx, c.TransferService.db.DBS().Reader); err != nil {
			return refs, err
		} else if vehicleHitLevel2Before {
			continue
		}

		firstTimeVIN, err := models.Vins(
			models.VinWhere.FirstEarningTokenID.EQ(types.NewDecimal(r.UserDeviceTokenID.Big)),
		).Exists(ctx, c.TransferService.db.DBS().Reader)
		if err != nil {
			return refs, err
		} else if !firstTimeVIN {
			continue
		}

		ownersOfLevel2FirstTimeVehicles[common.HexToAddress(r.UserEthereumAddress.String)] = struct{}{}
		numVehiclesLevel2FirstTime++
	}

	logger.Info().Msgf("Had %d VINs hit level 2 for the first time, with %d owners.", numVehiclesLevel2FirstTime, len(ownersOfLevel2FirstTimeVehicles))

	for user := range ownersOfLevel2FirstTimeVehicles {
		if userLevel2Before, err := models.Rewards(
			models.RewardWhere.IssuanceWeekID.LT(issuanceWeek),
			models.RewardWhere.ConnectionStreak.EQ(level2Weeks),
			models.RewardWhere.UserEthereumAddress.EQ(null.StringFrom(user.Hex())),
		).Exists(ctx, c.TransferService.db.DBS().Reader); err != nil {
			return refs, err
		} else if userLevel2Before {
			continue
		}

		if referredBefore, err := models.Referrals(models.ReferralWhere.Referee.EQ(user.Bytes())).Exists(ctx, c.TransferService.db.DBS().Reader); err != nil {
			return refs, err
		} else if referredBefore {
			continue
		}

		resp, err := c.UsersClient.GetUsersByEthereumAddress(ctx, &pb.GetUsersByEthereumAddressRequest{EthereumAddress: user.Bytes()})
		if err != nil {
			return refs, err
		}

		for _, potUser := range resp.Users {
			if potUser.ReferredBy == nil {
				continue
			}

			referrerID := "DELETED"
			referrerAddr := c.ContractAddress

			if potUser.ReferredBy.ReferrerValid && common.BytesToAddress(potUser.ReferredBy.EthereumAddress) != user {
				referrerID = potUser.ReferredBy.Id
				referrerAddr = common.BytesToAddress(potUser.ReferredBy.EthereumAddress)
			}

			refs.RefereeUserIDs = append(refs.RefereeUserIDs, potUser.Id)
			refs.Referees = append(refs.Referees, user)
			refs.ReferrerUserIDs = append(refs.ReferrerUserIDs, referrerID)
			refs.Referrers = append(refs.Referrers, referrerAddr)

			break
		}
	}

	logger.Info().Msgf("Sending out %d referrals.", numVehiclesLevel2FirstTime, len(ownersOfLevel2FirstTimeVehicles))

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

			refereesBatch := refs.Referees[i:j]
			referrersBatch := refs.Referrers[i:j]
			refereeIDsBatch := refs.RefereeUserIDs[i:j]
			referrerIDsBatch := refs.ReferrerUserIDs[i:j]

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

			for n := range refereesBatch {
				referrerID := null.StringFrom(referrerIDsBatch[n])
				if referrerID.String == "" {
					referrerID.Valid = false
				}
				r := models.Referral{
					Referee:        refereesBatch[n].Bytes(),
					Referrer:       referrersBatch[n].Bytes(),
					RequestID:      reqID,
					IssuanceWeekID: c.Week,
					RefereeUserID:  refereeIDsBatch[n],
					ReferrerUserID: referrerID,
				}
				if err := r.Insert(ctx, tx, boil.Infer()); err != nil {
					return err
				}
			}

			if err := c.BatchTransferReferralBonuses(reqID, refereesBatch, referrersBatch); err != nil {
				return err
			}

			return tx.Commit()
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
