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
	// MerkleDistributorAddress is the address of the MerkleDistributor contract.
	MerkleDistributorAddress string `yaml:"MERKLE_DISTRIBUTOR_ADDRESS"`
	// MerklePoolID is the id of the pool on the MerkleDistributor contract used
	// for baseline rewards. Defaults to 0.
	MerklePoolID int `yaml:"MERKLE_POOL_ID"`
	// FirstMerkleWeek is the first issuance week distributed via Merkle claims
	// instead of push transfers. A value of 0 disables the Merkle path.
	FirstMerkleWeek int `yaml:"FIRST_MERKLE_WEEK"`
	// MerkleTreeS3Bucket is the S3 bucket to which weekly Merkle tree files are uploaded.
	MerkleTreeS3Bucket string `yaml:"MERKLE_TREE_S3_BUCKET"`
	// MerkleTreeBaseURI is the public base URI under which uploaded tree files
	// are served, e.g. https://merkle.dimo.zone. The proofs URI passed to
	// setRoot is MERKLE_TREE_BASE_URI/pool-{poolId}/week-{week}.json.
	MerkleTreeBaseURI string `yaml:"MERKLE_TREE_BASE_URI"`
}
