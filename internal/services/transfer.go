package services

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/IBM/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/segmentio/ksuid"
)

type TransferService struct {
	Producer     sarama.SyncProducer
	RequestTopic string
	db           db.Store
	batchSize    int
}

func NewTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	db db.Store) *TransferService {

	return &TransferService{
		Producer:     producer,
		RequestTopic: settings.MetaTransactionSendTopic,
		db:           db,
		batchSize:    settings.TransferBatchSize,
	}
}

func (ts *TransferService) sendRequest(requestID string, addr common.Address, data []byte) error {
	event := shared.CloudEvent[transferData]{
		ID:          ksuid.New().String(),
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     requestID,
		Time:        time.Now(),
		Type:        "zone.dimo.transaction.request",
		Data: transferData{
			ID:   requestID,
			To:   addr,
			Data: data,
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = ts.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: ts.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	return err
}

type transferData struct {
	ID   string         `json:"id"`
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}

func (c *BaselineClient) transferTokens(ctx context.Context) error {
	batchSize := c.TransferService.batchSize
	responseSize := batchSize

	// If responseSize < pageSize then there must be no more pages of unsubmitted rewards.
	for batchSize == responseSize {
		reqID := ksuid.New().String()
		metaTxRequest := &models.MetaTransactionRequest{
			ID:     reqID,
			Status: models.MetaTransactionRequestStatusUnsubmitted,
		}

		err := metaTxRequest.Insert(ctx, c.TransferService.db.DBS().Writer, boil.Infer())
		if err != nil {
			return err
		}

		transfer, err := models.Rewards(
			qm.Expr(
				qm.Or2(models.RewardWhere.SyntheticDeviceTokens.GT(types.NewNullDecimal(decimal.New(0, 0)))),
				qm.Or2(models.RewardWhere.AftermarketDeviceTokens.GT(types.NewNullDecimal(decimal.New(0, 0)))),
				qm.Or2(models.RewardWhere.StreakTokens.GT(types.NewNullDecimal(decimal.New(0, 0)))),
			),
			models.RewardWhere.IssuanceWeekID.EQ(c.Week),
			models.RewardWhere.TransferMetaTransactionRequestID.IsNull(),
			// Temporary blacklist, see PLA-765.
			qm.LeftOuterJoin("rewards_api."+models.TableNames.Blacklist+" ON "+models.BlacklistTableColumns.UserEthereumAddress+" = "+models.RewardTableColumns.UserEthereumAddress),
			qm.Where(models.BlacklistTableColumns.UserEthereumAddress+" IS NULL"),
			qm.Limit(batchSize),
		).All(ctx, c.TransferService.db.DBS().Reader)
		if err != nil {
			return err
		}

		responseSize = len(transfer)

		tx, err := c.TransferService.db.DBS().Writer.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		transfers := []contracts.RewardTransferInfo{}

		for _, row := range transfer {
			trx := contracts.RewardTransferInfo{
				User:                       common.HexToAddress(row.RewardsReceiverEthereumAddress.String),
				VehicleId:                  row.UserDeviceTokenID.Int(nil),
				AftermarketDeviceId:        utils.NullDecimalToIntDefaultZero(row.AftermarketTokenID),
				ValueFromAftermarketDevice: utils.NullDecimalToIntDefaultZero(row.AftermarketDeviceTokens),
				SyntheticDeviceId:          big.NewInt(int64(row.SyntheticDeviceID.Int)),
				ValueFromSyntheticDevice:   utils.NullDecimalToIntDefaultZero(row.SyntheticDeviceTokens),
				ConnectionStreak:           big.NewInt(int64(row.ConnectionStreak)),
				ValueFromStreak:            utils.NullDecimalToIntDefaultZero(row.StreakTokens),
			}

			transfers = append(transfers, trx)
			row.TransferMetaTransactionRequestID = null.StringFrom(reqID)

			_, err = row.Update(ctx, tx, boil.Whitelist(models.RewardColumns.TransferMetaTransactionRequestID))
			if err != nil {
				return err
			}
		}

		err = c.BatchTransfer(reqID, transfers)
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BaselineClient) BatchTransfer(requestID string, transferInfo []contracts.RewardTransferInfo) error {
	abi, err := contracts.RewardMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("batchTransfer", transferInfo)
	if err != nil {
		return err
	}
	return c.TransferService.sendRequest(requestID, c.ContractAddress, data)
}
