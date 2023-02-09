package services

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

const (
	TransferEvent = "Transfer"
)

type kafkaEventStreamConsumer struct {
	Db  db.Store
	log *zerolog.Logger
}

type contractEventData struct {
	Contract        string          `json:"contract,omitempty"`
	TransactionHash string          `json:"transactionHash,omitempty"`
	Arguments       json.RawMessage `json:"arguments,omitempty"`
	Index           uint            `json:"index,omitempty"`
	BlockCompleted  bool            `json:"blockCompleted,omitempty"`
	EventSignature  string          `json:"eventSignature,omitempty"`
	EventName       string          `json:"eventName,omitempty"`
}

type transferEventData struct {
	Value *big.Int `json:"value"`
	To    string   `json:"to"`
	From  string   `json:"from"`
}

func NewEventConsumer(db db.Store, logger *zerolog.Logger) (*kafkaEventStreamConsumer, error) {
	return &kafkaEventStreamConsumer{Db: db, log: logger}, nil
}

func (c *kafkaEventStreamConsumer) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *kafkaEventStreamConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *kafkaEventStreamConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			event := shared.CloudEvent[contractEventData]{}
			err := json.Unmarshal(message.Value, &event)
			if err != nil {
				c.log.Err(err).Msg("error unmarshaling event")
				continue
			}

			switch event.Data.EventName {
			case TransferEvent:
				err = c.processTransferEvent(&event)
				if err != nil {
					c.log.Err(err).Msg("error processing transfer event")
				}
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (ec *kafkaEventStreamConsumer) processTransferEvent(e *shared.CloudEvent[contractEventData]) error {

	args := transferEventData{}
	err := json.Unmarshal(e.Data.Arguments, &args)

	transfer := models.TokenTransfer{
		AddressFrom:     args.From,
		AddressTo:       args.To,
		Amount:          types.NewDecimal(new(decimal.Big).SetBigMantScale(args.Value, 0)),
		TransactionHash: []byte(e.Subject),
		LogIndex:        int(e.Data.Index),
		BlockTimestamp:  e.Time,
	}

	err = transfer.Insert(context.Background(), ec.Db.DBS().Writer, boil.Infer())
	if err != nil {
		ec.log.Error().Err(err).Msg("Failed to insert token transfer record.")
		return err
	}

	return nil
}
