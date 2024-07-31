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

// CollectReferrals returns address pairs for referrals completed in the given week.
// These will come from referees who are earning for the first time and have a referrer
// attached to their account.
func (c *ReferralsClient) CollectReferrals(ctx context.Context, issuanceWeek int) (Referrals, error) {
	var refs Referrals

	covered := make(map[common.Address]struct{})

	maybeTriggering, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
		models.RewardWhere.ConnectionStreak.EQ(level2Weeks),
		models.RewardWhere.DisconnectionStreak.EQ(0),
	).All(ctx, c.TransferService.db.DBS().Reader)
	if err != nil {
		return Referrals{}, nil
	}

	for _, r := range maybeTriggering {
		user, err := c.UsersClient.GetUser(ctx, &pb.GetUserRequest{Id: r.UserID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				continue
			}
			return refs, err
		}

		if user.ReferredBy == nil {
			continue
		}

		if user.EthereumAddress == nil {
			c.Logger.Warn().Str("userId", r.UserID).Msg("Referred user does not have a valid ethereum address.")
			continue
		}

		refereeAddr := common.HexToAddress(*user.EthereumAddress)

		if _, ok := covered[refereeAddr]; ok {
			continue
		}

		referredBefore, err := models.Referrals(
			models.ReferralWhere.Referee.EQ(refereeAddr.Bytes()),
		).Exists(ctx, c.TransferService.db.DBS().Reader)
		if err != nil {
			return refs, err
		} else if referredBefore {
			covered[refereeAddr] = struct{}{}
			continue
		}

		if !user.ReferredBy.ReferrerValid {
			c.Logger.Info().Str("userId", r.UserID).Msg("Referring user has deleted their account or no longer has a confirmed ethereum address.")
			// referring eth addr is set to the referrals contract
			user.ReferredBy.EthereumAddress = c.ContractAddress.Bytes()
		}

		userHitBefore, err := models.Rewards(
			models.RewardWhere.UserEthereumAddress.EQ(null.StringFrom(*user.EthereumAddress)),
			models.RewardWhere.ConnectionStreak.EQ(level2Weeks),
			models.RewardWhere.IssuanceWeekID.LT(issuanceWeek),
		).Exists(ctx, c.TransferService.db.DBS().Reader)
		if err != nil {
			return refs, err
		} else if userHitBefore {
			covered[refereeAddr] = struct{}{}
			continue
		}

		first, err := models.Vins(
			models.VinWhere.FirstEarningTokenID.EQ(types.NewDecimal(r.UserDeviceTokenID.Big)),
		).Exists(ctx, c.TransferService.db.DBS().Reader)
		if err != nil {
			return refs, err
		} else if !first {
			continue
		}

		carHitLevel2Before, err := models.Rewards(
			models.RewardWhere.IssuanceWeekID.LT(issuanceWeek),
			models.RewardWhere.ConnectionStreak.EQ(level2Weeks),
			models.RewardWhere.UserDeviceTokenID.EQ(r.UserDeviceTokenID),
		).Exists(ctx, c.TransferService.db.DBS().Reader)
		if err != nil {
			return refs, err
		} else if carHitLevel2Before {
			continue
		}

		referrerAddr := common.BytesToAddress(user.ReferredBy.EthereumAddress)

		if refereeAddr == referrerAddr {
			c.Logger.Warn().Str("userId", r.UserID).Msg("Referred user's ethereum address is same as referring user's.")
			continue
		}

		covered[refereeAddr] = struct{}{}

		refs.Referees = append(refs.Referees, refereeAddr)
		refs.Referrers = append(refs.Referrers, referrerAddr)
		refs.RefereeUserIDs = append(refs.RefereeUserIDs, user.Id)
		refs.ReferrerUserIDs = append(refs.ReferrerUserIDs, user.ReferredBy.Id)
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
