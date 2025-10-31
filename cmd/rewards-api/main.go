package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	_ "github.com/DIMO-Network/rewards-api/docs"
	"github.com/DIMO-Network/rewards-api/internal/api"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/controllers"
	pb_tesla "github.com/DIMO-Network/tesla-oracle/pkg/grpc"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/services/identity"
	"github.com/DIMO-Network/rewards-api/internal/services/mobileapi"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	pb_rewards "github.com/DIMO-Network/shared/api/rewards"
	"github.com/DIMO-Network/shared/pkg/db"
	"github.com/DIMO-Network/shared/pkg/settings"

	"github.com/IBM/sarama"
	"github.com/burdiyan/kafkautil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var settingsPath = "settings.yaml"

// @title                       DIMO Rewards API
// @version                     1.0
// @BasePath                    /v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", "rewards-api").Logger()

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) == 40 {
				logger = logger.With().Str("commit", s.Value[:7]).Logger()
				break
			}
		}
	}

	ctx := context.Background()
	settings, err := settings.LoadConfig[config.Settings](settingsPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load settings.")
	}
	if settings.LogLevel != "" {
		lvl, err := zerolog.ParseLevel(settings.LogLevel)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to parse log level.")
		}
		zerolog.SetGlobalLevel(lvl)
	}

	if len(os.Args) == 1 {
		pdb := db.NewDbConnectionFromSettings(ctx, &settings.DB, true)
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return ErrorHandler(c, err, &logger, settings.Environment == "prod")
			},
		})

		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		devicesConn, err := grpc.NewClient(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices API client.")
		}
		defer devicesConn.Close()

		definitionsConn, err := grpc.NewClient(settings.DefinitionsAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create device definitions API client.")
		}
		defer definitionsConn.Close()

		definitionsClient := pb_defs.NewDeviceDefinitionServiceClient(definitionsConn)
		deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)
		chClient, err := ch.NewClient(&settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create ClickHouse client.")
		}

		rewardsController := controllers.RewardsController{
			DB:                pdb,
			Logger:            &logger,
			DefinitionsClient: definitionsClient,
			DevicesClient:     deviceClient,
			ChClient:          chClient,
			Settings:          &settings,
		}

		deviceController := controllers.DeviceController{
			DB:     pdb,
			Logger: &logger,
		}

		// secured paths
		keyRefreshInterval := time.Hour
		keyRefreshUnknownKID := true
		jwtAuth := jwtware.New(jwtware.Config{
			KeySetURL:            settings.JWTKeySetURL,
			KeyRefreshInterval:   &keyRefreshInterval,
			KeyRefreshUnknownKID: &keyRefreshUnknownKID,
		})
		app.Get("/v1/swagger/*", swagger.HandlerDefault)

		v1 := app.Group("/v1")

		v1.Get("/aftermarket/device/:tokenID", deviceController.GetDevice)
		v1.Get("/rewards/convert", rewardsController.GetHistoricalConversion)
		user := v1.Group("/user", jwtAuth)
		user.Get("/", rewardsController.GetUserRewards)
		user.Get("/history", rewardsController.GetUserRewardsHistory)
		user.Get("/history/transactions", rewardsController.GetTransactionHistory)
		user.Get("/history/balance", rewardsController.GetBalanceHistory)

		go startGRPCServer(&settings, pdb, &logger)

		// Start metatransaction listener
		kclient, err := createKafkaClient(&settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
		}

		consumer, err := sarama.NewConsumerGroupFromClient(settings.ConsumerGroup, kclient)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to initialize consumer group.")
		}
		defer consumer.Close()

		// need to pass logger here
		statusProc, err := services.NewStatusProcessor(pdb, &logger, &settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create transaction status processor.")
		}

		go func() {
			err := services.Consume(ctx, consumer, &settings, statusProc)
			if err != nil {
				logger.Fatal().Err(err).Send()
			}
		}()

		logger.Info().Msgf("Starting HTTP server on port %s.", settings.Port)
		if err := app.Listen(":" + settings.Port); err != nil {
			logger.Fatal().Err(err).Msgf("Fiber server failed.")
		}
		return
	}

	switch subCommand := os.Args[1]; subCommand {
	case "migrate":
		command := "up"
		if len(os.Args) > 2 {
			command = os.Args[2]
			if command == "down-to" || command == "up-to" {
				command = command + " " + os.Args[3]
			}
		}
		if err := database.MigrateDatabase(logger, &settings.DB, command, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run migration.")
		}
	case "calculate":
		var week int
		if len(os.Args) == 2 {
			// We have to subtract 1 because we're getting the number of the newly beginning week.
			week = date.GetWeekNumForCron(time.Now()) - 1
		} else {
			var err error
			week, err = strconv.Atoi(os.Args[2])
			if err != nil {
				logger.Fatal().Err(err).Msg("Could not parse week number.")
			}
		}
		pdb := db.NewDbConnectionFromSettings(ctx, &settings.DB, true)
		totalTime := 0
		for !pdb.IsReady() {
			if totalTime > 30 {
				logger.Fatal().Msg("could not connect to postgres after 30 seconds")
			}
			time.Sleep(time.Second)
			totalTime++
		}

		kclient, err := createKafkaClient(&settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
		}

		producer, err := sarama.NewSyncProducerFromClient(kclient)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka producer.")
		}

		transferService := services.NewTokenTransferService(&settings, producer, pdb)

		devicesConn, err := grpc.NewClient(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices-api connection.")
		}
		defer devicesConn.Close()

		deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)

		chClient, err := ch.NewClient(&settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create ClickHouse client.")
		}

		identClient := &identity.Client{
			QueryURL: settings.IdentityQueryURL,
			Client:   &http.Client{},
		}

		teslaConn, err := grpc.NewClient(settings.TeslaOracleGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices-api connection.")
		}
		defer teslaConn.Close()

		teslaClient := pb_tesla.NewTeslaOracleClient(teslaConn)

		baselineRewardClient := services.NewBaselineRewardService(&settings, transferService, chClient, deviceClient, identClient, week, &logger, teslaClient)

		if err := baselineRewardClient.BaselineIssuance(); err != nil {
			logger.Fatal().Err(err).Int("issuanceWeek", week).Msg("Failed to calculate and/or transfer rewards.")
		}
	case "issue-referral-bonus":
		var week int
		if len(os.Args) == 2 {
			// We have to subtract 1 because we're getting the number of the newly beginning week.
			week = date.GetWeekNumForCron(time.Now()) - 1
		} else {
			var err error
			week, err = strconv.Atoi(os.Args[2])
			if err != nil {
				logger.Fatal().Err(err).Msg("Could not parse week number.")
			}
		}

		logger := logger.With().Int("week", week).Logger()
		addr := common.HexToAddress(settings.ReferralContractAddress)
		logger.Info().Msgf("Running referral job with address %s.", addr)

		pdb := db.NewDbConnectionFromSettings(ctx, &settings.DB, true)
		totalTime := 0
		for !pdb.IsReady() {
			if totalTime > 30 {
				logger.Fatal().Msg("could not connect to postgres after 30 seconds")
			}
			time.Sleep(time.Second)
			totalTime++
		}

		kclient, err := createKafkaClient(&settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
		}

		producer, err := sarama.NewSyncProducerFromClient(kclient)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka producer.")
		}

		transferService := services.NewTokenTransferService(&settings, producer, pdb)

		mapURL, err := url.ParseRequestURI(settings.MobileAPIBaseURL)
		if err != nil {
			logger.Fatal().Err(err).Msg("Couldn't parse Mobile API URL.")
		}

		mAPI := mobileapi.New(mapURL)

		referralsClient := services.NewReferralBonusService(&settings, transferService, week, &logger, mAPI)
		err = referralsClient.ReferralsIssuance(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to transfer referral bonuses.")
		}
	case "describe":
		tokenID, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.Fatal().Err(err).Msg("Couldn't parse vehicle token id.")
		}

		identClient := &identity.Client{
			QueryURL: settings.IdentityQueryURL,
			Client:   &http.Client{},
		}

		vehDesc, err := identClient.DescribeVehicle(uint64(tokenID))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to describe vehicle.")
		}

		b, err := json.MarshalIndent(vehDesc, "", "  ")
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to marshal vehicle description.")
		}

		fmt.Println(string(b))
	case "owner":
		ok := common.IsHexAddress(os.Args[2])
		if !ok {
			logger.Fatal().Msg("Couldn't parse owner address.")
		}

		owner := common.HexToAddress(os.Args[2])

		identClient := &identity.Client{
			QueryURL: settings.IdentityQueryURL,
			Client:   &http.Client{},
		}

		vehDesc, err := identClient.GetVehicles(owner)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to describe vehicle.")
		}

		b, err := json.MarshalIndent(vehDesc, "", "  ")
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to marshal vehicle description.")
		}

		fmt.Println(string(b))
	default:
		logger.Fatal().Msgf("Unrecognized sub-command %s.", subCommand)
	}
}

func startGRPCServer(settings *config.Settings, dbs db.Store, logger *zerolog.Logger) {
	lis, err := net.Listen("tcp", ":"+settings.GRPCPort)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Couldn't listen on gRPC port %s", settings.GRPCPort)
	}

	logger.Info().Msgf("Starting gRPC server on port %s", settings.GRPCPort)
	server := grpc.NewServer()
	pb_rewards.RegisterRewardsServiceServer(server, api.NewRewardsService(dbs, logger))

	if err := server.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("gRPC server terminated unexpectedly")
	}
}

// logging stuff here

func GetLogger(c *fiber.Ctx, d *zerolog.Logger) *zerolog.Logger {
	m := c.Locals("logger")
	if m == nil {
		return d
	}

	l, ok := m.(*zerolog.Logger)
	if !ok {
		return d
	}

	return l
}

// ErrorHandler custom handler to log recovered errors using our logger and return json instead of string
func ErrorHandler(c *fiber.Ctx, err error, logger *zerolog.Logger, isProduction bool) error {
	logger = GetLogger(c, logger)

	code := fiber.StatusInternalServerError // Default 500 statuscode

	e, fiberTypeErr := err.(*fiber.Error)
	if fiberTypeErr {
		// Override status code if fiber.Error type
		code = e.Code
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	codeStr := strconv.Itoa(code)

	logger.Err(err).Str("httpStatusCode", codeStr).
		Str("httpMethod", c.Method()).
		Str("httpPath", c.Path()).
		Msg("caught an error from http request")
	// return an opaque error if we're in a higher level environment and we haven't specified an fiber type err.
	if !fiberTypeErr && isProduction {
		err = fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	return c.Status(code).JSON(ErrorRes{
		Code:    code,
		Message: err.Error(),
	})
}

type ErrorRes struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func createKafkaClient(settings *config.Settings) (sarama.Client, error) {
	kconf := sarama.NewConfig()
	kconf.Version = sarama.V2_8_1_0
	kconf.Producer.Return.Successes = true
	kconf.Producer.Partitioner = kafkautil.NewJVMCompatiblePartitioner

	return sarama.NewClient(strings.Split(settings.KafkaBrokers, ","), kconf)
}
