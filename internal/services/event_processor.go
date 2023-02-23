package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

const (
	contractEventType = "zone.dimo.contract.event"
)

type ContractEventStreamConsumer struct {
	Db     db.Store
	log    *zerolog.Logger
	tokens map[string]string
}

type contractEventData struct {
	Contract        common.Address  `json:"contract,omitempty"`
	TransactionHash common.Hash     `json:"transactionHash,omitempty"`
	Arguments       json.RawMessage `json:"arguments,omitempty"`
	Block           eventBlock      `json:"block"`
	Index           int             `json:"index,omitempty"`
	EventSignature  common.Hash     `json:"eventSignature,omitempty"`
	EventName       string          `json:"eventName,omitempty"`
}

type eventBlock struct {
	Number int64
	Hash   common.Hash
	Time   time.Time
}

type TokenConfig struct {
	Tokens []struct {
		ChainID int64          `yaml:"chainId"`
		Address common.Address `yaml:"address"`
	} `yaml:"tokens"`
}

func NewEventConsumer(db db.Store, logger *zerolog.Logger, tc *TokenConfig) (*ContractEventStreamConsumer, error) {
	m := map[string]string{}

	for _, tk := range tc.Tokens {
		logger.Info().Msgf("Tracking %s on chain %d.", tk.Address, tk.ChainID)
		m[fmt.Sprintf("chain/%d", tk.ChainID)] = hexutil.Encode(tk.Address.Bytes())
	}

	return &ContractEventStreamConsumer{Db: db,
		log:    logger,
		tokens: m,
	}, nil
}

func (c *ContractEventStreamConsumer) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *ContractEventStreamConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *ContractEventStreamConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			var event shared.CloudEvent[contractEventData]
			err := json.Unmarshal(message.Value, &event)
			if err != nil {
				c.log.Err(err).Msg("error unmarshaling event")
				session.MarkMessage(message, "")
				continue
			}

			if event.Type != contractEventType {
				session.MarkMessage(message, "")
				continue
			}

			if addr, ok := c.tokens[event.Source]; !ok || addr != event.Subject || event.Data.EventName != "Transfer" {
				session.MarkMessage(message, "")
				continue
			}

			err = c.processTransferEvent(&event)
			if err != nil {
				c.log.Err(err).Str("contract", event.Data.Contract.Hex()).Str("txHash", event.Data.TransactionHash.Hex()).Int("logIndex", event.Data.Index).Msg("error storing transfer event")
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (ec *ContractEventStreamConsumer) processTransferEvent(e *shared.CloudEvent[contractEventData]) error {
	if !strings.HasPrefix(e.Source, "chain/") {
		return fmt.Errorf("source doesn't have the chain/ prefix: %s", e.Source)
	}

	chainIDRaw := strings.TrimPrefix(e.Source, "chain/")
	chainID, err := strconv.ParseInt(chainIDRaw, 10, 64)
	if err != nil {
		return fmt.Errorf("couldn't parse chain id %q: %w", chainIDRaw, err)
	}

	var args contracts.TokenTransfer
	err = json.Unmarshal(e.Data.Arguments, &args)
	if err != nil {
		ec.log.Error().Err(err).Msg("failed to unpack event arguments")
		return err
	}

	transfer := models.TokenTransfer{
		AddressFrom:     args.From.Bytes(),
		AddressTo:       args.To.Bytes(),
		Amount:          types.NewDecimal(new(decimal.Big).SetBigMantScale(args.Value, 0)),
		TransactionHash: e.Data.TransactionHash.Bytes(),
		LogIndex:        e.Data.Index,
		BlockTimestamp:  e.Data.Block.Time,
		ChainID:         chainID,
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
