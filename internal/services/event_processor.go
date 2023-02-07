package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/Shopify/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

const (
	TransferEvent = "Transfer"
)

type event struct {
	ID          string            `json:"id"`
	Source      string            `json:"source"`
	SpecVersion string            `json:"specversion"`
	Subject     string            `json:"subject"`
	Time        time.Time         `json:"time"`
	Type        string            `json:"type"`
	Data        contractEventData `json:"data"`
}

type eventConsumer struct {
	Db  func() *database.DBReaderWriter
	log *zerolog.Logger
}

type contractEventData struct {
	Contract        string                 `json:"contract,omitempty"`
	TransactionHash string                 `json:"transactionHash,omitempty"`
	Arguments       map[string]interface{} `json:"arguments,omitempty"`
	BlockCompleted  bool                   `json:"blockCompleted,omitempty"`
	EventSignature  string                 `json:"eventSignature,omitempty"`
	EventName       string                 `json:"eventName,omitempty"`
}

type transferEventData struct {
	Value float64 `json:"value"`
	To    string  `json:"to"`
	From  string  `json:"from"`
}

func NewEventConsumer(db func() *database.DBReaderWriter, logger *zerolog.Logger) (*eventConsumer, error) {
	return &eventConsumer{Db: db, log: logger}, nil
}

func (c *eventConsumer) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *eventConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *eventConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			event := event{}
			err := json.Unmarshal(message.Value, &event)
			if err != nil {
				return err
			}

			switch event.Data.EventName {
			case TransferEvent:
				err = c.processTransferEvent(&event)
				if err != nil {
					return err
				}
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (ec *eventConsumer) processTransferEvent(e *event) error {

	args := transferEventData{}
	b, err := json.Marshal(e.Data.Arguments)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &args)

	amnt := new(decimal.Big).SetFloat64(args.Value)
	amount := types.NewDecimal(amnt)

	transfer := models.TokenTransfer{
		ContractAddress: common.FromHex(e.Data.Contract),
		UserAddressFrom: common.FromHex(args.From),
		UserAddressTo:   common.FromHex(args.To),
		Amount:          amount,
		CreatedAt:       e.Time,
	}

	err = transfer.Insert(context.Background(), ec.Db().Writer, boil.Infer())
	if err != nil {
		ec.log.Error().Err(err).Msg("Failed to insert token transfer record.")
		return err
	}

	return nil
}
