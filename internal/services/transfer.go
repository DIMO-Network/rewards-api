package services

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type TransferService struct {
	Producer        sarama.SyncProducer
	Consumer        sarama.ConsumerGroup
	RequestTopic    string
	StatusTopic     string
	db              db.Store
	ContractAddress common.Address
	batchSize       int
}

func NewTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	// settings config.Settings, producer sarama.SyncProducer, reqTopic string, contract Contract,
	db db.Store) *TransferService {

	return &TransferService{
		Producer:     producer,
		RequestTopic: settings.MetaTransactionSendTopic,
		StatusTopic:  settings.MetaTransactionStatusTopic,
		db:           db,
		batchSize:    settings.TransferBatchSize,
	}
}

func (ts *TransferService) sendRequest(requestID string, addr common.Address, data []byte) error {
	event := shared.CloudEvent[transferData]{
		ID:          requestID,
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     addr.String(),
		Time:        time.Now(),
		Type:        "zone.dimo.transaction.request",
		Data: transferData{
			ID:   requestID,
			To:   hexutil.Encode(ts.ContractAddress[:]),
			Data: hexutil.Encode(data),
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
	ID   string `json:"id"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func (c *BaselineClient) transfer(ctx context.Context) error {
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
			models.RewardWhere.Tokens.GT(types.NewNullDecimal(decimal.New(0, 0))),
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

		userAddr := make([]common.Address, responseSize)
		tknValues := make([]*big.Int, responseSize)
		vehicleIds := make([]*big.Int, responseSize)

		tx, err := c.TransferService.db.DBS().Writer.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		for i, row := range transfer {
			userAddr[i] = common.HexToAddress(row.UserEthereumAddress.String)
			tknValues[i] = row.Tokens.Int(nil)
			vehicleIds[i] = row.UserDeviceTokenID.Int(nil)

			row.TransferMetaTransactionRequestID = null.StringFrom(reqID)

			_, err = row.Update(ctx, tx, boil.Whitelist(models.RewardColumns.TransferMetaTransactionRequestID))
			if err != nil {
				return err
			}
		}

		err = c.BatchTransfer(reqID, userAddr, tknValues, vehicleIds)
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

func (c *BaselineClient) BatchTransfer(requestID string, users []common.Address, values []*big.Int, vehicleIds []*big.Int) error {
	abi, err := contracts.RewardMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("batchTransfer", users, values, vehicleIds)
	if err != nil {
		return err
	}
	return c.TransferService.sendRequest(requestID, c.ContractAddress, data)
}
