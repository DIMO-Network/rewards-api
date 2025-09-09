package config

import (
	"github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
	"github.com/DIMO-Network/shared/pkg/db"
	"github.com/ethereum/go-ethereum/common"
)

// Settings contains the application config
type Settings struct {
	Environment                string          `yaml:"ENVIRONMENT"`
	Port                       string          `yaml:"PORT"`
	LogLevel                   string          `yaml:"LOG_LEVEL"`
	DB                         db.Settings     `yaml:"DB"`
	JWTKeySetURL               string          `yaml:"JWT_KEY_SET_URL"`
	GRPCPort                   string          `yaml:"GRPC_PORT"`
	DevicesAPIGRPCAddr         string          `yaml:"DEVICES_API_GRPC_ADDR"`
	DefinitionsAPIGRPCAddr     string          `yaml:"DEFINITIONS_API_GRPC_ADDR"`
	FetchAPIGRPCAddr           string          `yaml:"FETCH_API_GRPC_ADDR"`
	KafkaBrokers               string          `yaml:"KAFKA_BROKERS"`
	MetaTransactionSendTopic   string          `yaml:"META_TRANSACTION_SEND_TOPIC"`
	MetaTransactionStatusTopic string          `yaml:"META_TRANSACTION_STATUS_TOPIC"`
	IssuanceContractAddress    string          `yaml:"ISSUANCE_CONTRACT_ADDRESS"`
	ReferralContractAddress    string          `yaml:"REFERRAL_CONTRACT_ADDRESS"`
	ConsumerGroup              string          `yaml:"CONSUMER_GROUP"`
	TransferBatchSize          int             `yaml:"TRANSFER_BATCH_SIZE"`
	FirstAutomatedWeek         int             `yaml:"FIRST_AUTOMATED_WEEK"`
	ContractEventTopic         string          `yaml:"CONTRACT_EVENT_TOPIC"`
	Clickhouse                 config.Settings `yaml:",inline"`
	IdentityQueryURL           string          `yaml:"IDENTITY_QUERY_URL"`
	EnableStaking              bool            `yaml:"ENABLE_STAKING"`
	DIMORegistryChainID        int             `yaml:"DIMO_REGISTRY_CHAIN_ID"`
	VehicleNFTAddress          common.Address  `yaml:"VEHICLE_NFT_ADDRESS"`
	VINVCDataVersion           string          `yaml:"VINVC_DATA_VERSION"`
	MobileAPIBaseURL           string          `yaml:"MOBILE_API_BASE_URL"`
	StorageNodeDevLicense      common.Address  `yaml:"STORAGE_NODE_DEV_LICENSE"`
	VINVCConcurrencyLimit      int             `yaml:"VINVC_CONCURRENCY_LIMIT"`
}
