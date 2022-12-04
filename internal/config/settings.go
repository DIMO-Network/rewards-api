package config

import "fmt"

// Settings contains the application config
type Settings struct {
	Environment                    string `yaml:"ENVIRONMENT"`
	Port                           string `yaml:"PORT"`
	LogLevel                       string `yaml:"LOG_LEVEL"`
	DBUser                         string `yaml:"DB_USER"`
	DBPassword                     string `yaml:"DB_PASSWORD"`
	DBPort                         string `yaml:"DB_PORT"`
	DBHost                         string `yaml:"DB_HOST"`
	DBName                         string `yaml:"DB_NAME"`
	DBMaxOpenConnections           int    `yaml:"DB_MAX_OPEN_CONNECTIONS"`
	DBMaxIdleConnections           int    `yaml:"DB_MAX_IDLE_CONNECTIONS"`
	ServiceName                    string `yaml:"SERVICE_NAME"`
	JWTKeySetURL                   string `yaml:"JWT_KEY_SET_URL"`
	GRPCPort                       string `yaml:"GRPC_PORT"`
	SwaggerBaseURL                 string `yaml:"SWAGGER_BASE_URL"`
	ElasticSearchAnalyticsHost     string `yaml:"ELASTIC_SEARCH_ANALYTICS_HOST"`
	ElasticSearchAnalyticsUsername string `yaml:"ELASTIC_SEARCH_ANALYTICS_USERNAME"`
	ElasticSearchAnalyticsPassword string `yaml:"ELASTIC_SEARCH_ANALYTICS_PASSWORD"`
	DeviceDataIndexName            string `yaml:"DEVICE_DATA_INDEX_NAME"`
	DevicesAPIGRPCAddr             string `yaml:"DEVICES_API_GRPC_ADDR"`
	DefinitionsAPIGRPCAddr         string `yaml:"DEFINITIONS_API_GRPC_ADDR"`
	KafkaBrokers                   string `yaml:"KAFKA_BROKERS"`
	MetaTransactionSendTopic       string `yaml:"META_TRANSACTION_SEND_TOPIC"`
	MetaTransactionStatusTopic     string `yaml:"META_TRANSACTION_STATUS_TOPIC"`
	IssuanceContractAddress        string `yaml:"ISSUANCE_CONTRACT_ADDRESS"`
	ConsumerGroup                  string `yaml:"CONSUMER_GROUP"`
	TransferBatchSize              int    `yaml:"TRANSFER_BATCH_SIZE"`
	FirstAutomatedWeek             int    `yaml:"FIRST_AUTOMATED_WEEK"`
}

// GetWriterDSN builds the connection string to the db writer - for now same as reader
func (app *Settings) GetWriterDSN(withSearchPath bool) string {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		app.DBUser,
		app.DBPassword,
		app.DBName,
		app.DBHost,
		app.DBPort,
	)
	if withSearchPath {
		dsn = fmt.Sprintf("%s search_path=%s", dsn, app.DBName) // assumption is schema has same name as dbname
	}
	return dsn
}
