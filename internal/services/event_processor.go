package services

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/DIMO-Network/rewards-api/internal/config"
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
	TransferEvent  = "Transfer"
	BlockProcessed = "zone.dimo.blockchain.block.processed"
)

type ContractEventStreamConsumer struct {
	Db            db.Store
	log           *zerolog.Logger
	TokenContract string
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

func NewEventConsumer(db db.Store, logger *zerolog.Logger, settings *config.Settings) (*ContractEventStreamConsumer, error) {
	return &ContractEventStreamConsumer{Db: db, log: logger, TokenContract: settings.TokenAddress}, nil
}

func (c *ContractEventStreamConsumer) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *ContractEventStreamConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *ContractEventStreamConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			event := shared.CloudEvent[contractEventData]{}
			err := json.Unmarshal(message.Value, &event)
			if err != nil {
				c.log.Err(err).Msg("error unmarshaling event")
				session.MarkMessage(message, "")
				continue
			}

			if event.Type == BlockProcessed {
				session.MarkMessage(message, "")
				continue
			}

			switch event.Data.Contract {
			case c.TokenContract:
				switch event.Data.EventName {
				case TransferEvent:
					err = c.processTransferEvent(&event)
					if err != nil {
						c.log.Err(err).Str("contract", event.Data.Contract).Str("txHash", event.Data.TransactionHash).Int("logIndex", int(event.Data.Index)).Msg("error storing transfer event")
					}
				}
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (ec *ContractEventStreamConsumer) processTransferEvent(e *shared.CloudEvent[contractEventData]) error {

	args := transferEventData{}
	err := json.Unmarshal(e.Data.Arguments, &args)
	if err != nil {
		ec.log.Error().Err(err).Msg("failed to unpack event arguments")
		return err
	}

	transfer := models.TokenTransfer{
		AddressFrom:     []byte(args.From),
		AddressTo:       []byte(args.To),
		Amount:          types.NewDecimal(new(decimal.Big).SetBigMantScale(args.Value, 0)),
		TransactionHash: []byte(e.Subject),
		LogIndex:        int(e.Data.Index),
		BlockTimestamp:  e.Time,
	}

	err = transfer.Upsert(context.Background(), ec.Db.DBS().Writer, true,
		[]string{models.TokenTransferColumns.TransactionHash, models.TokenTransferColumns.LogIndex},
		boil.Infer(), boil.Infer())
	if err != nil {
		ec.log.Error().Err(err).Msg("failed to insert token transfer record.")
		return err
	}

	return nil
}
