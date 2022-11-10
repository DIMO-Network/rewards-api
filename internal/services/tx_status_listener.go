package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/big"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/Shopify/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
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

func NewStatusProcessor(pdb func() *database.DBReaderWriter, logger *zerolog.Logger) (*TransferStatusProcessor, error) {
	abi, err := contracts.RewardMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return &TransferStatusProcessor{
		ABI:                    abi,
		DB:                     pdb,
		Logger:                 logger,
		DidntQualifyEvent:      abi.Events["DidntQualify"],
		TokensTransferredEvent: abi.Events["TokensTransferred"],
	}, nil
}

type TransferStatusProcessor struct {
	ABI                    *abi.ABI
	DB                     func() *database.DBReaderWriter
	Logger                 *zerolog.Logger
	DidntQualifyEvent      abi.Event
	TokensTransferredEvent abi.Event
}

func (s *TransferStatusProcessor) processMessage(msg *sarama.ConsumerMessage) error {
	event := shared.CloudEvent[ceData]{}
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return err
	}

	s.Logger.Info().
		Interface("eventData", event.Data).
		Msg("Processing transaction status.")

	tx, err := s.DB().Writer.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() //nolint

	txnRow, err := models.FindMetaTransactionRequest(context.Background(), s.DB().Reader, event.Data.RequestID)
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

				var userDeviceTokenID *big.Int

				switch log.Topics[0] {
				case s.TokensTransferredEvent.ID:
					success = true
					event := contracts.RewardTokensTransferred{}
					err := s.parseLog(&event, s.TokensTransferredEvent, *txLog)
					if err != nil {
						return err
					}
					userDeviceTokenID = event.VehicleNodeId
				case s.DidntQualifyEvent.ID:
					event := contracts.RewardDidntQualify{}
					err := s.parseLog(&event, s.DidntQualifyEvent, *txLog)
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

func (s *TransferStatusProcessor) parseLog(out any, event abi.Event, log eth_types.Log) error {
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

func convertLog(logIn *ceLog) *eth_types.Log {
	return &eth_types.Log{
		Topics: logIn.Topics,
		Data:   logIn.Data,
	}
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
