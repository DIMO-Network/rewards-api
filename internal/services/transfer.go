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
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
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
		batchSize:       settings.TransferBatchSize,
	}
}

type Client struct {
	Producer        sarama.SyncProducer
	Consumer        sarama.ConsumerGroup
	RequestTopic    string
	StatusTopic     string
	db              func() *database.DBReaderWriter
	ContractAddress common.Address
	batchSize       int
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
	batchSize := c.batchSize
	responseSize := batchSize

	// If responseSize < pageSize then there must be no more pages of unsubmitted rewards.
	for batchSize == responseSize {
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
			models.RewardWhere.Tokens.GT(types.NewNullDecimal(decimal.New(0, 0))),
			models.RewardWhere.IssuanceWeekID.EQ(week),
			models.RewardWhere.TransferMetaTransactionRequestID.IsNull(),
			// Temporary blacklist, see PLA-765.
			models.RewardWhere.UserEthereumAddress.NEQ(null.StringFrom("0x481e8DB1dd18fd02caA8A83Ef7A73cF207b83930")),
			models.RewardWhere.UserEthereumAddress.NEQ(null.StringFrom("0x3596Da3ab608d4fD63F4Bc9F4631A6838d435474")),
			models.RewardWhere.UserEthereumAddress.NEQ(null.StringFrom("0x21762721Fe155F29D2EdbBB2a88688a032c41c58")),
			models.RewardWhere.UserEthereumAddress.NEQ(null.StringFrom("0xBE421ef2988794F8061A62FE8A45BA29e08458C6")),
			models.RewardWhere.UserEthereumAddress.NEQ(null.StringFrom("0x620a84B6D68a33017109c6cBECa233442b89c237")),
			models.RewardWhere.UserEthereumAddress.NEQ(null.StringFrom("0x7088A745eED70B7348678577095eA332d4f9A3Dd")),
			qm.Limit(batchSize),
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
