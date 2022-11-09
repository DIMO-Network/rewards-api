package services

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Transfer interface {
	// BatchTransfer(requestID string, users []common.Address, values []*big.Int, vehicleIds []string) error
	TransferUserTokens(ctx context.Context, week int) error
}

func NewTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	contractAddress common.Address,
	// settings config.Settings, producer sarama.SyncProducer, reqTopic string, contract Contract,
	db func() *database.DBReaderWriter) Transfer {

	return &Client{
		ContractAddress: contractAddress,
		Producer:        producer,
		RequestTopic:    settings.MetaTransactionSendTopic,
		StatusTopic:     settings.MetaTransactionStatusTopic,
		db:              db,
	}
}

type Client struct {
	Producer        sarama.SyncProducer
	Consumer        sarama.ConsumerGroup
	RequestTopic    string
	StatusTopic     string
	db              func() *database.DBReaderWriter
	ContractAddress common.Address
	Settings        *config.Settings
}

type transferData struct {
	ID   string `json:"id"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func (c *Client) TransferUserTokens(ctx context.Context, week int) error {
	err := c.transfer(ctx, week)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) transfer(ctx context.Context, week int) error {
	pageSize := 2
	responseSize := pageSize

	// If responseSize < pageSize then there must be no more pages of unsubmitted rewards.
	for pageSize == responseSize {
		reqID := ksuid.New().String()
		metaTxRequest := &models.MetaTransactionRequest{
			ID:     reqID,
			Status: models.MetaTransactionRequestStatusUnsubmitted,
		}

		err := metaTxRequest.Insert(ctx, c.db().Writer, boil.Infer())
		if err != nil {
			return err
		}

		transfer, err := models.Rewards(
			models.RewardWhere.IssuanceWeekID.EQ(week),
			models.RewardWhere.TransferMetaTransactionRequestID.IsNull(),
			qm.Limit(pageSize),
		).All(ctx, c.db().Reader)
		if err != nil {
			return err
		}

		responseSize = len(transfer)
		userAddr := make([]common.Address, responseSize)
		tknValues := make([]*big.Int, responseSize)
		vehicleIds := make([]*big.Int, responseSize)

		tx, err := c.db().Writer.BeginTx(ctx, nil)
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

func (c *Client) BatchTransfer(requestID string, users []common.Address, values []*big.Int, vehicleIds []*big.Int) error {
	abi, err := contracts.RewardMetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack("batchTransfer", users, values, vehicleIds)
	if err != nil {
		return err
	}
	return c.sendRequest(requestID, data)
}

func (c *Client) sendRequest(requestID string, data []byte) error {
	event := shared.CloudEvent[transferData]{
		ID:          ksuid.New().String(),
		Source:      "rewards-api",
		SpecVersion: "1.0",
		Subject:     requestID,
		Time:        time.Now(),
		Type:        "zone.dimo.transaction.request",
		Data: transferData{
			ID:   requestID,
			To:   hexutil.Encode(c.ContractAddress[:]),
			Data: hexutil.Encode(data),
		},
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = c.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: c.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	return err
}
