package services

import (
	"context"
	"encoding/json"

	"github.com/DIMO-Network/shared"
	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
)

var requestsTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "reward_issuance_processor",
		Subsystem: "consumer",
		Name:      "requests_total",
	},
)

type consumer struct {
	logger *zerolog.Logger
	client *ethclient.Client
}

type TransferEventData struct {
	ID   string `json:"id"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		requestsTotal.Inc()

		event := &shared.CloudEvent[TransferEventData]{}

		err := json.Unmarshal(msg.Value, event)
		if err != nil {
			c.logger.Err(err).Int32("partition", msg.Partition).Int64("offset", msg.Offset).Msg("Couldn't parse message, skipping.")
			continue
		}

		c.logger.Info().Interface("user", event.ID).Msg("Token transfer request recieved.")

		// to := common.HexToAddress(event.Data.To)
		// tkns := common.FromHex(event.Data.Data)

		c.logger.Info().Interface("tokens", event.Data).Msg("Token amount.")
		session.MarkMessage(msg, "")
	}
	return nil
}

func New(ctx context.Context, name string, topic string, kafkaClient sarama.Client, logger *zerolog.Logger, ethClient *ethclient.Client) error {
	group, err := sarama.NewConsumerGroupFromClient(name, kafkaClient)
	if err != nil {
		return err
	}

	consumer := &consumer{logger: logger, client: ethClient} // manager: manager

	for {
		err := group.Consume(ctx, []string{topic}, consumer)
		if err != nil {
			logger.Err(err).Msg("Consumer group session did not terminate gracefully.")
		}
		if ctx.Err() != nil {
			// Context canceled, so quit.
			return nil
		}
	}
}
