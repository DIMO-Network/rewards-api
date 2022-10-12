package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	issuance "github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Transfer interface {
	BatchTransfer(requestID string, devices []deviceInfo) error
	TransferUserTokens(week int, ctx context.Context) error
}

func NewTokenTransferService(
	// settings config.Settings, producer sarama.SyncProducer, reqTopic string, contract Contract,
	db func() *database.DBReaderWriter) Transfer {
	return &Client{
		// Producer:     producer,
		// RequestTopic: reqTopic,
		// Contract: Contract{
		// 	ChainID: settings.ChainID,
		// 	Address: settings.Address,
		// 	Name:    settings.ContractName,
		// 	Version: settings.ContractVersion},
		db: db}
}

type Client struct {
	Producer     sarama.SyncProducer
	RequestTopic string
	Contract     Contract
	db           func() *database.DBReaderWriter
}

type Contract struct {
	ChainID *big.Int
	Address common.Address
	Name    string
	Version string
}

type transferData struct {
	ID   string `json:"id"`
	To   string `json:"to"`
	Data string `json:"data"`
}

type deviceInfo struct {
	To          common.Address
	Amount      *big.Int
	VehicleNode string
}

func (c *Client) TransferUserTokens(week int, ctx context.Context) error {

	offset := 0
	batchSize := 2
	responseLength := 101

	for batchSize <= responseLength {
		rewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(week), qm.Limit(batchSize), qm.Offset(offset*batchSize)).All(ctx, c.db().Reader)
		if err != nil {
			return err
		}

		devices := make([]deviceInfo, len(rewards))

		for n, row := range rewards {
			devices[n].Amount = row.Tokens.Int(nil)
			devices[n].VehicleNode = row.UserDeviceID // confirm the type that this should be, uint?
			// fetch user address from users api
			devices[n].To = common.Address{}
		}

		reqID := fmt.Sprintf("%d-Request %d", week, offset+1)
		err = c.BatchTransfer(reqID, devices)
		if err != nil {
			return err
		}

		offset++
		responseLength = len(rewards)
	}

	return nil

}

func (c *Client) BatchTransfer(requestID string, devices []deviceInfo) error {
	abi, err := issuance.IssuanceMetaData.GetAbi()
	if err != nil {
		return err
	}

	data, err := abi.Pack("batchTransfer", devices)
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
			To:   hexutil.Encode(c.Contract.Address[:]),
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
