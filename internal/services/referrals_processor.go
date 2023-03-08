package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type referralsConsumerGroupHandler struct {
	rsp *ReferralProcessor
}

func (referralsConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (referralsConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h referralsConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := h.rsp.processMessage(msg)
		if err != nil {
			return err
		}
		sess.MarkMessage(msg, "")
	}

	return nil
}

func ConsumeReferrals(ctx context.Context, group sarama.ConsumerGroup, settings *config.Settings, refProc *ReferralProcessor) error {
	for {
		topics := []string{settings.MetaTransactionStatusTopic}
		handler := referralsConsumerGroupHandler{rsp: refProc}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			return err
		}
	}
}

func NewReferralStatusProcessor(pdb db.Store, logger *zerolog.Logger) (*ReferralProcessor, error) {
	abi, err := contracts.ReferralsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &ReferralProcessor{
		ABI:              abi,
		DB:               pdb,
		Logger:           logger,
		ReferralComplete: abi.Events["ReferralComplete"],
		ReferralInvalid:  abi.Events["ReferralInvalid"],
	}, nil
}

type ReferralProcessor struct {
	ABI              *abi.ABI
	DB               db.Store
	Logger           *zerolog.Logger
	ReferralComplete abi.Event
	ReferralInvalid  abi.Event
}

func (s *ReferralProcessor) processMessage(msg *sarama.ConsumerMessage) error {
	event := shared.CloudEvent[ceData]{}
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return err
	}

	s.Logger.Info().
		Interface("eventData", event.Data).
		Msg("Processing transaction status.")

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

	txnRow.Hash = null.StringFrom(event.Data.Transaction.Hash)
	txnRow.Status = event.Data.Type

	if event.Data.Type == "Confirmed" {
		txnRow.Successful = null.BoolFrom(*event.Data.Transaction.Successful)

		if *event.Data.Transaction.Successful {
			for _, log := range event.Data.Transaction.Logs {
				success := false
				txLog := convertLog(&log)

				var user common.Address
				var referrer common.Address
				switch log.Topics[0] {
				case s.ReferralComplete.ID:
					success = true
					rc := contracts.ReferralsReferralComplete{}
					err := s.parseLog(&event, s.ReferralComplete, *txLog)
					if err != nil {
						return err
					}
					user = rc.Referred
					referrer = rc.Referrer

				case s.ReferralInvalid.ID:
					ri := contracts.ReferralsReferralInvalid{}
					err := s.parseLog(&event, s.ReferralInvalid, *txLog)
					if err != nil {
						return err
					}
					user = ri.Referred
					referrer = ri.Referrer

				default:
					continue
				}

				refRow, err := models.Referrals(
					models.ReferralWhere.ID.EQ(event.Data.RequestID),
					models.ReferralWhere.Referred.EQ(user[:]),
					models.ReferralWhere.Referrer.EQ(referrer[:]),
				).One(context.Background(), tx)
				if err != nil {
					return err
				}

				refRow.JobStatus = models.ReferralsJobStatusComplete
				refRow.TransferSuccessful = null.BoolFrom(success)
				if !success {
					refRow.TransferFailureReason = null.StringFrom(models.ReferralsTransferFailureReasonReferralInvalid)
				}

				_, err = refRow.Update(context.TODO(), tx, boil.Whitelist(models.ReferralColumns.TransferSuccessful, models.ReferralColumns.TransferFailureReason))
				if err != nil {
					return err
				}
			}
		} else {
			_, err := models.Referrals(
				models.ReferralWhere.ID.EQ(event.Data.RequestID),
			).UpdateAll(context.Background(), tx, models.M{
				models.ReferralColumns.TransferSuccessful:    false,
				models.ReferralColumns.TransferFailureReason: models.ReferralsTransferFailureReasonTxReverted,
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

func (s *ReferralProcessor) parseLog(out any, event abi.Event, log eth_types.Log) error {
	if len(log.Data) > 0 {
		err := s.ABI.UnpackIntoInterface(out, event.Name, log.Data)
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
