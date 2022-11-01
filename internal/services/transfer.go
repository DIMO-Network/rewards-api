package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	issuance "github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Transfer interface {
	// BatchTransfer(requestID string, users []common.Address, values []*big.Int, vehicleIds []string) error
	TransferUserTokens(ctx context.Context, week int) error
}

type consumerGroupHandler struct {
	name string
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	abi, err := issuance.IssuanceMetaData.GetAbi()
	if err != nil {
		return err
	}

	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		return err
	}
	pdb := database.NewDbConnectionFromSettings(context.Background(), &settings)

	temp := &S{ABI: abi, DB: pdb.DBS}
	for msg := range claim.Messages() {
		err := temp.processMessages(msg)
		if err != nil {
			return err
		}
		sess.MarkMessage(msg, "")
	}

	return nil
}

func Consume(group sarama.ConsumerGroup, wg *sync.WaitGroup, name string) {
	defer wg.Done()
	ctx := context.Background()
	for {
		topics := []string{"topic.transaction.request.status"}
		handler := consumerGroupHandler{name: name}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func NewTokenTransferService(
	settings *config.Settings,
	producer sarama.SyncProducer,
	usersClient pb_users.UserServiceClient,
	devicesClient pb_devices.UserDeviceServiceClient,
	contractAddress common.Address,
	// settings config.Settings, producer sarama.SyncProducer, reqTopic string, contract Contract,
	db func() *database.DBReaderWriter) Transfer {

	return &Client{
		usersClient:     usersClient,
		devicesClient:   devicesClient,
		ContractAddress: contractAddress,
		Producer:        producer,
		RequestTopic:    "topic.transaction.request.send",
		StatusTopic:     "topic.transaction.request.status",
		db:              db}
}

type Client struct {
	Producer        sarama.SyncProducer
	Consumer        sarama.ConsumerGroup
	RequestTopic    string
	StatusTopic     string
	db              func() *database.DBReaderWriter
	usersClient     pb_users.UserServiceClient
	devicesClient   pb_devices.UserDeviceServiceClient
	ContractAddress common.Address
}

type Contract struct {
	Address common.Address
	Name    string
	Version string
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

	page := 0
	pageSize := 2
	responseSize := pageSize

	for pageSize == responseSize {
		reqID := ksuid.New().String()

		q := `
		INSERT INTO rewards_api.meta_transaction_requests
		(id, status) VALUES ($1, $2);`
		_, err := c.db().Writer.ExecContext(ctx, q, reqID, "Unsubmitted")
		if err != nil {
			return err
		}

		transfer, err := models.Rewards(
			models.RewardWhere.IssuanceWeekID.EQ(week),
			qm.OrderBy(models.RewardColumns.UserDeviceID), // Somewhat dangerous, what if this changes?
			qm.Limit(pageSize),
			qm.Offset(page*pageSize),
		).All(ctx, c.db().Reader)
		if err != nil {
			return err
		}

		responseSize = len(transfer)
		var userAddr []common.Address
		var tknValues []*big.Int
		var vehicleIds []*big.Int

		tx, _ := c.db().GetWriterConn().BeginTx(ctx, nil)
		stmt, _ := tx.Prepare(`UPDATE rewards_api.rewards
			SET transfer_meta_transaction_request_id = $1
			WHERE user_id = $2 AND user_device_id = $3;`)
		for _, row := range transfer {
			userAddr = append(userAddr, common.HexToAddress(row.UserEthereumAddress.String))
			tknValues = append(tknValues, row.Tokens.Int(nil))
			vehicleIds = append(vehicleIds, row.UserDeviceTokenID.Int(nil))
			_, err := stmt.Exec(reqID, row.UserID, row.UserDeviceID)
			if err != nil {
				rollbackErr := tx.Rollback()

				if rollbackErr != nil {
					combinedErrors := fmt.Sprintf("error rolling back transaction: %s", rollbackErr.Error())
					err = errors.Wrap(err, combinedErrors)
				}
				return err
			}

		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		err = c.BatchTransfer(reqID, userAddr, tknValues, vehicleIds)
		if err != nil {
			return err
		}

		page++
	}
	return nil
}

func (c *Client) BatchTransfer(requestID string, users []common.Address, values []*big.Int, vehicleIds []*big.Int) error {
	abi, err := issuance.IssuanceMetaData.GetAbi()
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

	log.Printf("topic=%q, requestId=%q", c.RequestTopic, requestID)

	p, o, err := c.Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: c.RequestTopic,
			Key:   sarama.StringEncoder(requestID),
			Value: sarama.ByteEncoder(eventBytes),
		},
	)

	log.Printf("err=%v, partition=%d, offset=%d", err, p, o)

	return err
}

func (s *S) processMessages(msg *sarama.ConsumerMessage) error {
	var DidNotQualifyHash = crypto.Keccak256Hash([]byte("DidntQualify(address,uint256,uint256)"))
	event := shared.CloudEvent[ceData]{}
	err := json.Unmarshal(msg.Value, &event)
	if err != nil {
		return err
	}
	abi, err := issuance.IssuanceMetaData.GetAbi()
	if err != nil {
		return err
	}

	tx, _ := s.DB().GetWriterConn().BeginTx(context.Background(), nil)
	stmt, _ := tx.Prepare(`UPDATE rewards_api.rewards
			SET transfer_successful = $1
			WHERE transfer_meta_transaction_request_id = $2 AND user_device_token_id = $3;`)

	logs := event.Data.Transaction.Logs
	for _, l := range logs {
		var success bool
		txLog := convertLog(&l)
		rec := issuance.IssuanceTokensTransferred{}
		err := s.parseLog(&rec, abi.Events["TokensTransferred"], *txLog)
		if err != nil {
			return err
		}

		if !(l.Topics[0] == DidNotQualifyHash.String()) {
			success = *event.Data.Transaction.Successful
			_, err = stmt.Exec(success, event.Data.RequestID, rec.VehicleNodeId.Int64())
		} else {
			_, err = stmt.Exec(success, event.Data.RequestID, rec.VehicleNodeId.Int64())
		}

		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				combinedErrors := fmt.Sprintf("error rolling back transaction: %s", rollbackErr.Error())
				err = errors.Wrap(err, combinedErrors)
			}
			return err
		}

	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	q := `
	UPDATE rewards_api.meta_transaction_requests
	SET hash = $1, status = $2, successful = $3
	WHERE id = $4;`
	_, err = s.DB().Writer.Exec(q, event.Data.Transaction.Hash, event.Data.Type, event.Data.Transaction.Successful, event.Data.RequestID)
	if err != nil {
		return err
	}

	return nil

}

type S struct {
	ABI    *abi.ABI
	DB     func() *database.DBReaderWriter
	Logger *zerolog.Logger
}

func (s *S) parseLog(out any, event abi.Event, log eth_types.Log) error {
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
	topics := make([]common.Hash, len(logIn.Topics))
	for i, t := range logIn.Topics {
		topics[i] = common.HexToHash(t)
	}

	data := common.FromHex(logIn.Data)

	return &eth_types.Log{
		Topics: topics,
		Data:   data,
	}
}

type ceLog struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
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
