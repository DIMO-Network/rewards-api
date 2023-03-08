package config

import (
	"github.com/DIMO-Network/shared/db"
)

// Settings contains the application config
type Settings struct {
	Environment                    string      `yaml:"ENVIRONMENT"`
	Port                           string      `yaml:"PORT"`
	LogLevel                       string      `yaml:"LOG_LEVEL"`
	DB                             db.Settings `yaml:"DB"`
	ServiceName                    string      `yaml:"SERVICE_NAME"`
	JWTKeySetURL                   string      `yaml:"JWT_KEY_SET_URL"`
	GRPCPort                       string      `yaml:"GRPC_PORT"`
	SwaggerBaseURL                 string      `yaml:"SWAGGER_BASE_URL"`
	ElasticSearchAnalyticsHost     string      `yaml:"ELASTIC_SEARCH_ANALYTICS_HOST"`
	ElasticSearchAnalyticsUsername string      `yaml:"ELASTIC_SEARCH_ANALYTICS_USERNAME"`
	ElasticSearchAnalyticsPassword string      `yaml:"ELASTIC_SEARCH_ANALYTICS_PASSWORD"`
	DeviceDataIndexName            string      `yaml:"DEVICE_DATA_INDEX_NAME"`
	DevicesAPIGRPCAddr             string      `yaml:"DEVICES_API_GRPC_ADDR"`
	DefinitionsAPIGRPCAddr         string      `yaml:"DEFINITIONS_API_GRPC_ADDR"`
	UsersAPIGRPCAddr               string      `yaml:"USERS_API_GRPC_ADDR"`
	KafkaBrokers                   string      `yaml:"KAFKA_BROKERS"`
	MetaTransactionSendTopic       string      `yaml:"META_TRANSACTION_SEND_TOPIC"`
	MetaTransactionStatusTopic     string      `yaml:"META_TRANSACTION_STATUS_TOPIC"`
	IssuanceContractAddress        string      `yaml:"ISSUANCE_CONTRACT_ADDRESS"`
	ConsumerGroup                  string      `yaml:"CONSUMER_GROUP"`
	TransferBatchSize              int         `yaml:"TRANSFER_BATCH_SIZE"`
	FirstAutomatedWeek             int         `yaml:"FIRST_AUTOMATED_WEEK"`
	ContractEventTopic             string      `yaml:"CONTRACT_EVENT_TOPIC"`
}
