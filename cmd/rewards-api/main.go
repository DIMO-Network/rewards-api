package main

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	_ "github.com/DIMO-Network/rewards-api/docs"
	"github.com/DIMO-Network/rewards-api/internal/api"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/DIMO-Network/rewards-api/internal/controllers"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/shared"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	pb_rewards "github.com/DIMO-Network/shared/api/rewards"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/burdiyan/kafkautil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// @title                       DIMO Rewards API
// @version                     1.0
// @BasePath                    /v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", "rewards-api").Logger()

	ctx := context.Background()
	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load settings.")
	}

	if len(os.Args) == 1 {
		pdb := db.NewDbConnectionFromSettings(ctx, &settings.DB, true)
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler:          ErrorHandler,
		})

		app.Get("/", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		devicesConn, err := grpc.Dial(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices API client.")
		}
		defer devicesConn.Close()

		definitionsConn, err := grpc.Dial(settings.DefinitionsAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create device definitions API client.")
		}
		defer definitionsConn.Close()

		usersConn, err := grpc.Dial(settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create device definitions API client.")
		}
		defer usersConn.Close()

		definitionsClient := pb_defs.NewDeviceDefinitionServiceClient(definitionsConn)
		deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)
		usersClient := pb_users.NewUserServiceClient(usersConn)

		dataClient := services.NewDeviceDataClient(&settings)

		f, err := os.Open("config.yaml")
		if err != nil {
			logger.Fatal().Err(err).Msg("Couldn't load config.")
		}
		defer f.Close()

		cs, err := io.ReadAll(f)
		if err != nil {
			logger.Fatal().Err(err).Msg("Couldn't load config file.")
		}

		rc := os.ExpandEnv(string(cs))

		var tc services.TokenConfig
		if err := yaml.Unmarshal([]byte(rc), &tc); err != nil {
			logger.Fatal().Err(err).Msg("Couldn't load token config.")
		}

		var tks []*contracts.Token

		for _, tkCopy := range tc.Tokens {
			client, err := ethclient.Dial(tkCopy.RPCURL)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Failed to create client for chain %d.", tkCopy.ChainID)
			}

			token, err := contracts.NewToken(tkCopy.Address, client)
			if err != nil {
				logger.Fatal().Err(err).Msgf("Failed to instantiate token for chain %d.", tkCopy.ChainID)
			}

			tks = append(tks, token)
		}

		rewardsController := controllers.RewardsController{
			DB:                pdb,
			Logger:            &logger,
			DefinitionsClient: definitionsClient,
			DevicesClient:     deviceClient,
			DataClient:        dataClient,
			Settings:          &settings,
			Tokens:            tks,
			UsersClient:       usersClient,
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

		v1 := app.Group("/v1", jwtAuth)
		v1.Get("/user", rewardsController.GetUserRewards)
		v1.Get("/user/history", rewardsController.GetUserRewardsHistory)
		v1.Get("/user/history/transactions", rewardsController.GetTransactionHistory)
		v1.Get("/user/history/balance", rewardsController.GetBalanceHistory)

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

		kclient2, err := createKafkaClient(&settings)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
		}

		consumer2, err := sarama.NewConsumerGroupFromClient(settings.ConsumerGroup, kclient2)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka consumer.")
		}

		msgHandler, err := services.NewEventConsumer(pdb, &logger, &tc)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create new event consumer.")
		}

		go func() {
			for {
				err = consumer2.Consume(ctx, []string{settings.ContractEventTopic}, msgHandler)
				if err != nil {
					logger.Err(err).Msg("error while processing messages")
					if ctx.Err() != nil {
						return
					}
				}
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
		logger.Info().Msg("XDD")
		if err := database.MigrateDatabase(logger, &settings.DB, command, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run migration.")
		}
	case "calculate":
		var week int
		if len(os.Args) == 2 {
			// We have to subtract 1 because we're getting the number of the newly beginning week.
			week = services.GetWeekNumForCron(time.Now()) - 1
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

		devicesConn, err := grpc.Dial(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices-api connection.")
		}
		defer devicesConn.Close()

		deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)

		definitionsConn, err := grpc.Dial(settings.DefinitionsAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create device-definitions-api connection.")
		}
		defer definitionsConn.Close()

		definitionsClient := pb_defs.NewDeviceDefinitionServiceClient(definitionsConn)

		baselineRewardClient := services.NewBaselineRewardService(&settings, transferService, services.NewDeviceDataClient(&settings), deviceClient, definitionsClient, week, &logger)

		if err := baselineRewardClient.Calculate(week); err != nil {
			logger.Fatal().Err(err).Int("issuanceWeek", week).Msg("Failed to calculate and/or transfer rewards.")
		}
	case "issue-referral-bonus":
		var week int
		if len(os.Args) == 2 {
			// We have to subtract 1 because we're getting the number of the newly beginning week.
			week = services.GetWeekNumForCron(time.Now()) - 1
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

		usersConn, err := grpc.Dial(settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create device definitions API client.")
		}
		defer usersConn.Close()

		usersClient := pb_users.NewUserServiceClient(usersConn)

		referralsClient := services.NewReferralBonusService(&settings, transferService, week, &logger, usersClient)
		err = referralsClient.ReferralsIssuance(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to transfer referral bonuses.")
		}

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

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError // Default.
	message := "Internal error."

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = err.Error()
	}

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": message,
	})
}

func createKafkaClient(settings *config.Settings) (sarama.Client, error) {
	kconf := sarama.NewConfig()
	kconf.Version = sarama.V2_8_1_0
	kconf.Producer.Return.Successes = true
	kconf.Producer.Partitioner = kafkautil.NewJVMCompatiblePartitioner

	return sarama.NewClient(strings.Split(settings.KafkaBrokers, ","), kconf)
}
