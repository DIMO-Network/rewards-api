package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/big"

	"github.com/DIMO-Network/cloudevent"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/pkg/db"
	"github.com/IBM/sarama"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type consumerGroupHandler struct {
	tsp *TransferStatusProcessor
}

func (consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := h.tsp.processMessage(msg)
		if err != nil {
			return err
		}
		sess.MarkMessage(msg, "")
	}

	return nil
}

func Consume(ctx context.Context, group sarama.ConsumerGroup, settings *config.Settings, statusProc *TransferStatusProcessor) error {
	for {
		topics := []string{settings.MetaTransactionStatusTopic}
		handler := consumerGroupHandler{tsp: statusProc}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			return err
		}
	}
}

func NewStatusProcessor(pdb db.Store, logger *zerolog.Logger, settings *config.Settings) (*TransferStatusProcessor, error) {
	baselineABI, err := contracts.RewardMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	referralABI, err := contracts.ReferralMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &TransferStatusProcessor{
		DB:     pdb,
		Logger: logger,
		BaselineProcessor: &BaselineProcessor{
			Address:                common.HexToAddress(settings.IssuanceContractAddress),
			ABI:                    baselineABI,
			DidntQualifyEvent:      baselineABI.Events["DidntQualify"],
			TokensTransferredEvent: baselineABI.Events["TokensTransferred"],
		},
		ReferralsProcessor: &ReferralsProcessor{
			Address:          common.HexToAddress(settings.ReferralContractAddress),
			ABI:              referralABI,
			ReferralInvalid:  referralABI.Events["ReferralInvalid"],
			ReferralComplete: referralABI.Events["ReferralComplete"],
		},
	}, nil
}

type TransferStatusProcessor struct {
	ReferralsProcessor *ReferralsProcessor
	BaselineProcessor  *BaselineProcessor
	DB                 db.Store
	Logger             *zerolog.Logger
}

type ReferralsProcessor struct {
	Address          common.Address
	ABI              *abi.ABI
	ReferralInvalid  abi.Event
	ReferralComplete abi.Event
}

type BaselineProcessor struct {
	Address                common.Address
	ABI                    *abi.ABI
	DidntQualifyEvent      abi.Event
	TokensTransferredEvent abi.Event
}

func (s *TransferStatusProcessor) processMessage(msg *sarama.ConsumerMessage) error {
	event := cloudevent.CloudEvent[ceData]{}
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return err
	}

	s.Logger.Info().Msg("Processing transaction status.")

	mtr, err := models.MetaTransactionRequests(
		models.MetaTransactionRequestWhere.ID.EQ(event.Data.RequestID),
		qm.Load(models.MetaTransactionRequestRels.TransferMetaTransactionRequestRewards),
		qm.Load(models.MetaTransactionRequestRels.RequestReferrals),
	).One(context.TODO(), s.DB.DBS().Reader)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	switch {
	case len(mtr.R.TransferMetaTransactionRequestRewards) != 0:
		err := s.processBaselineEvent(event)
		if err != nil {
			return err
		}
	case len(mtr.R.RequestReferrals) != 0:
		err := s.processReferralEvent(event)
		if err != nil {
			return err
		}
	default:
		s.Logger.Error().Msgf("Known meta-transaction %s has no associated baseline or referral batch.", event.Data.RequestID)
	}

	return nil
}

func (s *TransferStatusProcessor) processBaselineEvent(event cloudevent.CloudEvent[ceData]) error {
	tx, err := s.DB.DBS().Writer.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint

	txnRow, err := models.FindMetaTransactionRequest(context.Background(), s.DB.DBS().Reader, event.Data.RequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	txnRow.Status = event.Data.Type

	if event.Data.Type != models.MetaTransactionRequestStatusFailed {
		txnRow.Hash = null.StringFrom(event.Data.Transaction.Hash)
	}

	if event.Data.Type == models.MetaTransactionRequestStatusConfirmed {
		txnRow.Successful = null.BoolFrom(*event.Data.Transaction.Successful)

		if *event.Data.Transaction.Successful {
			for _, log := range event.Data.Transaction.Logs {
				success := false

				var userDeviceTokenID *big.Int

				switch log.Topics[0] {
				case s.BaselineProcessor.TokensTransferredEvent.ID:
					success = true
					event := contracts.RewardTokensTransferred{}
					err := s.parseLog(&event, s.BaselineProcessor.TokensTransferredEvent, log, s.BaselineProcessor.ABI)
					if err != nil {
						return err
					}
					userDeviceTokenID = event.VehicleNodeId
				case s.BaselineProcessor.DidntQualifyEvent.ID:
					event := contracts.RewardDidntQualify{}
					err := s.parseLog(&event, s.BaselineProcessor.DidntQualifyEvent, log, s.BaselineProcessor.ABI)
					if err != nil {
						return err
					}

					userDeviceTokenID = event.VehicleNodeId
				default:
					continue
				}

				rewardRow, err := models.Rewards(
					models.RewardWhere.TransferMetaTransactionRequestID.EQ(null.StringFrom(event.Data.RequestID)),
					models.RewardWhere.UserDeviceTokenID.EQ(types.NewNullDecimal(new(decimal.Big).SetBigMantScale(userDeviceTokenID, 0))),
				).One(context.Background(), tx)
				if err != nil {
					return err
				}

				rewardRow.TransferSuccessful = null.BoolFrom(success)
				if !success {
					rewardRow.TransferFailureReason = null.StringFrom(models.RewardsTransferFailureReasonDidntQualify)
				}

				_, err = rewardRow.Update(context.TODO(), tx, boil.Whitelist(models.RewardColumns.TransferSuccessful, models.RewardColumns.TransferFailureReason))
				if err != nil {
					return err
				}
			}
		} else {
			_, err := models.Rewards(
				models.RewardWhere.TransferMetaTransactionRequestID.EQ(null.StringFrom(event.Data.RequestID)),
			).UpdateAll(context.Background(), tx, models.M{
				models.RewardColumns.TransferSuccessful:    false,
				models.RewardColumns.TransferFailureReason: models.RewardsTransferFailureReasonTxReverted,
			})
			if err != nil {
				return err
			}
		}
	}

	_, err = txnRow.Update(context.TODO(), tx, boil.Infer())
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *TransferStatusProcessor) processReferralEvent(cloudEvent cloudevent.CloudEvent[ceData]) error {
	tx, err := s.DB.DBS().Writer.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint

	txnRow, err := models.FindMetaTransactionRequest(context.Background(), s.DB.DBS().Reader, cloudEvent.Data.RequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// TODO(elffjs): This is probably a very bad scenario.
			return nil
		}
		return err
	}

	txnRow.Status = cloudEvent.Data.Type

	if cloudEvent.Data.Type != models.MetaTransactionRequestStatusFailed {
		txnRow.Hash = null.StringFrom(cloudEvent.Data.Transaction.Hash)
	}

	if cloudEvent.Data.Type == models.MetaTransactionRequestStatusConfirmed {
		txnRow.Successful = null.BoolFrom(*cloudEvent.Data.Transaction.Successful)

		if *cloudEvent.Data.Transaction.Successful {
			for _, log := range cloudEvent.Data.Transaction.Logs {
				success := false

				var referee, referrer common.Address

				switch log.Topics[0] {
				case s.ReferralsProcessor.ReferralComplete.ID:
					success = true
					var event contracts.ReferralReferralComplete
					err := s.parseLog(&event, s.ReferralsProcessor.ReferralComplete, log, s.ReferralsProcessor.ABI)
					if err != nil {
						return err
					}
					referee, referrer = event.Referee, event.Referrer
				case s.ReferralsProcessor.ReferralInvalid.ID:
					var event contracts.ReferralReferralInvalid
					err := s.parseLog(&event, s.ReferralsProcessor.ReferralInvalid, log, s.ReferralsProcessor.ABI)
					if err != nil {
						return err
					}
					referee, referrer = event.Referee, event.Referrer

				default:
					continue
				}

				rewardRow, err := models.Referrals(
					models.ReferralWhere.Referee.EQ(referee.Bytes()),
					models.ReferralWhere.Referrer.EQ(referrer.Bytes()),
				).One(context.Background(), tx)
				if err != nil {
					return err
				}

				rewardRow.TransferSuccessful = null.BoolFrom(success)
				if !success {
					rewardRow.TransferFailureReason = null.StringFrom(models.ReferralsTransferFailureReasonReferralInvalid)
				}

				_, err = rewardRow.Update(context.TODO(), tx, boil.Whitelist(models.ReferralColumns.TransferSuccessful, models.RewardColumns.TransferFailureReason))
				if err != nil {
					return err
				}
			}
		} else {
			if _, err := models.Referrals(
				models.ReferralWhere.RequestID.EQ(cloudEvent.Data.RequestID),
			).UpdateAll(context.Background(), tx, models.M{
				models.ReferralColumns.TransferSuccessful:    false,
				models.ReferralColumns.TransferFailureReason: models.ReferralsTransferFailureReasonTxReverted,
			}); err != nil {
				return err
			}
		}
	}

	_, err = txnRow.Update(context.TODO(), tx, boil.Infer())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *TransferStatusProcessor) parseLog(out any, event abi.Event, log ceLog, ctrABI *abi.ABI) error {
	if len(log.Data) > 0 {
		err := ctrABI.UnpackIntoInterface(out, event.Name, log.Data)
		if err != nil {
			return err
		}
	}

	var indexed abi.Arguments
	for _, arg := range event.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}

	err := abi.ParseTopics(out, indexed, log.Topics[1:])
	if err != nil {
		return err
	}

	return nil
}

type ceLog struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    hexutil.Bytes  `json:"data"`
}

type ceTx struct {
	Hash       string  `json:"hash"`
	Successful *bool   `json:"successful,omitempty"`
	Logs       []ceLog `json:"logs,omitempty"`
}

// Just using the same struct for all three event types. Lazy.
type ceData struct {
	RequestID   string `json:"requestId"`
	Type        string `json:"type"`
	Transaction ceTx   `json:"transaction"`
}
