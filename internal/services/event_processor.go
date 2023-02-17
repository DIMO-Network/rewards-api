package services

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

const (
	BlockProcessed = "zone.dimo.blockchain.block.processed"
)

type ContractEventStreamConsumer struct {
	Db                 db.Store
	log                *zerolog.Logger
	TokenAddress       string
	TokenTransferEvent abi.Event
}

type contractEventData struct {
	Contract        string          `json:"contract,omitempty"`
	TransactionHash string          `json:"transactionHash,omitempty"`
	Arguments       json.RawMessage `json:"arguments,omitempty"`
	Block           block           `json:"block"`
	Index           uint            `json:"index,omitempty"`
	BlockCompleted  bool            `json:"blockCompleted,omitempty"`
	EventSignature  string          `json:"eventSignature,omitempty"`
	EventName       string          `json:"eventName,omitempty"`
}

type block struct {
	Number int64
	Hash   string
	Time   time.Time
}

type transferEventData struct {
	Value *big.Int       `json:"value"`
	To    common.Address `json:"to"`
	From  common.Address `json:"from"`
}

type Config struct {
	TokenAddress string `yaml:"token_address"`
	ChainID      string `yaml:"chain_id"`
}

func NewEventConsumer(db db.Store, logger *zerolog.Logger, conf *Config) (*ContractEventStreamConsumer, error) {
	abi, err := contracts.TokenMetaData.GetAbi()
	if err != nil {
		return &ContractEventStreamConsumer{}, err
	}

	return &ContractEventStreamConsumer{Db: db,
		log:                logger,
		TokenAddress:       conf.TokenAddress,
		TokenTransferEvent: abi.Events["Transfer"]}, nil
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
			case c.TokenAddress:
				switch event.Data.EventName {
				case c.TokenTransferEvent.Name:
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
		AddressFrom:     args.From[:],
		AddressTo:       args.To[:],
		Amount:          types.NewDecimal(new(decimal.Big).SetBigMantScale(args.Value, 0)),
		TransactionHash: []byte(e.Subject),
		LogIndex:        int(e.Data.Index),
		BlockTimestamp:  e.Time,
		ChainID:         e.Source,
	}

	err = transfer.Upsert(context.Background(), ec.Db.DBS().Writer, true,
		[]string{models.TokenTransferColumns.TransactionHash, models.TokenTransferColumns.LogIndex, models.TokenTransferColumns.ChainID},
		boil.Infer(), boil.Infer())
	if err != nil {
		ec.log.Error().Err(err).Msg("failed to insert token transfer record.")
		return err
	}

	return nil
}
