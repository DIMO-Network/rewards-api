package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/IBM/sarama"
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

type DIMOTransferListener struct {
	Db     db.Store
	log    *zerolog.Logger
	tokens map[string]string

	// Only one partition for contract events, so no need to share the cache.
	// But we lock anyway just so it doesn't look weird.
	userAddrsmap map[common.Address]struct{}

	outsideAddrsmap map[common.Address]struct{}

	mu sync.Mutex
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
		RPCURL  string         `yaml:"rpcUrl"`
	} `yaml:"tokens"`
}

func NewEventConsumer(db db.Store, logger *zerolog.Logger, tc *TokenConfig) (*DIMOTransferListener, error) {
	m := map[string]string{}

	for _, tk := range tc.Tokens {
		logger.Info().Msgf("Tracking %s on chain %d.", tk.Address, tk.ChainID)
		m[fmt.Sprintf("chain/%d", tk.ChainID)] = hexutil.Encode(tk.Address.Bytes())
	}

	return &DIMOTransferListener{Db: db,
		log:    logger,
		tokens: m,
	}, nil
}

func (c *DIMOTransferListener) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *DIMOTransferListener) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *DIMOTransferListener) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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

			// The source field is, e.g., "chain/137" for Polygon.
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

func (c *DIMOTransferListener) areTransferAddrsRelevant(from, to common.Address) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.userAddrsmap[from]; ok {
		return true, nil
	}
	if _, ok := c.userAddrsmap[to]; ok {
		return true, nil
	}

	// Neither address known to be a user.

	if _, ok := c.outsideAddrsmap[from]; ok {
		if _, ok := c.outsideAddrsmap[to]; ok {
			return false, nil
		}

	}

	// TODO(elffjs): Could try to record who's not a user
}

func (c *DIMOTransferListener) processTransferEvent(e *shared.CloudEvent[contractEventData]) error {
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
		c.log.Error().Err(err).Msg("failed to unpack event arguments")
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

	err = transfer.Upsert(context.Background(), c.Db.DBS().Writer, true,
		[]string{models.TokenTransferColumns.TransactionHash, models.TokenTransferColumns.LogIndex, models.TokenTransferColumns.ChainID},
		boil.Infer(), boil.Infer())
	if err != nil {
		c.log.Error().Err(err).Msg("failed to insert token transfer record.")
		return err
	}

	return nil
}
