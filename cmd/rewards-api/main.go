package main

import (
	"context"
	"errors"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	_ "github.com/DIMO-Network/rewards-api/docs"
	"github.com/DIMO-Network/rewards-api/internal/api"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/controllers"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/shared"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	pb_rewards "github.com/DIMO-Network/shared/api/rewards"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/Shopify/sarama"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/burdiyan/kafkautil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var baseWeeklyTokens = new(big.Int).Mul(big.NewInt(1_105_000), ether)

// @title                       DIMO Rewards API
// @version                     1.0
// @BasePath                    /v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	ctx := context.Background()
	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		os.Exit(1)
	}

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "rewards-api").
		Logger()

	if len(os.Args) == 1 {
		pdb := database.NewDbConnectionFromSettings(ctx, &settings)
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
		defer devicesConn.Close()

		definitionsClient := pb_defs.NewDeviceDefinitionServiceClient(definitionsConn)
		deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)

		dataClient := services.NewDeviceDataClient(&settings)

		rewardsController := controllers.RewardsController{
			DB:                pdb.DBS,
			Logger:            &logger,
			DefinitionsClient: definitionsClient,
			DevicesClient:     deviceClient,
			DataClient:        dataClient,
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

		if settings.Environment != "prod" {
			v1.Get("/points", rewardsController.GetPointsThisWeek)
			v1.Get("/tokens", rewardsController.GetTokensThisWeek)
			v1.Get("/allocation", rewardsController.GetUserAllocation)
		}

		go startGRPCServer(&settings, pdb.DBS, &logger)

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
		migrateDatabase(logger, &settings, command, "rewards_api")
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
		pdb := database.NewDbConnectionFromSettings(ctx, &settings)
		totalTime := 0
		for !pdb.IsReady() {
			if totalTime > 30 {
				logger.Fatal().Msg("could not connect to postgres after 30 seconds")
			}
			time.Sleep(time.Second)
			totalTime++
		}
		task := services.RewardsTask{
			Settings:    &settings,
			DataService: services.NewDeviceDataClient(&settings),
			DB:          pdb.DBS,
			Logger:      &logger,
		}
		if err := task.Calculate(week); err != nil {
			logger.Fatal().Err(err).Int("issuanceWeek", week).Msg("Failed to calculate rewards.")
		}
	case "test-tokens":
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
		pdb := database.NewDbConnectionFromSettings(ctx, &settings)
		totalTime := 0
		for !pdb.IsReady() {
			if totalTime > 30 {
				logger.Fatal().Msg("could not connect to postgres after 30 seconds")
			}
			time.Sleep(time.Second)
			totalTime++
		}

		st := storage.NewDB(pdb.DBS)
		err := st.AssignTokens(ctx, week, baseWeeklyTokens)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to assign tokens.")
		}
		rewards, err := st.Rewards(ctx, week)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to retrieve rewards.")
		}
		logger.Info().Interface("rewards", rewards).Msg("Output.")
	case "test-transfer":
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
		pdb := database.NewDbConnectionFromSettings(ctx, &settings)
		totalTime := 0
		for !pdb.IsReady() {
			if totalTime > 30 {
				logger.Fatal().Msg("could not connect to postgres after 30 seconds")
			}
			time.Sleep(time.Second)
			totalTime++
		}
		devicesConn, err := grpc.Dial(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices API client.")
		}
		defer devicesConn.Close()

		devicesClient := pb_devices.NewUserDeviceServiceClient(devicesConn)
		usersConn, err := grpc.Dial(settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices API client.")
		}
		defer usersConn.Close()
		usersClient := pb_users.NewUserServiceClient(usersConn)

		kconf := sarama.NewConfig()
		kconf.Version = sarama.V2_8_1_0
		kconf.Producer.Return.Successes = true
		kconf.Producer.Partitioner = kafkautil.NewJVMCompatiblePartitioner

		kclient, err := sarama.NewClient(strings.Split(settings.KafkaBrokers, ","), kconf)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
		}

		log.Printf("brokers=%q", settings.KafkaBrokers)

		producer, err := sarama.NewSyncProducerFromClient(kclient)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka producer.")
		}
		addr := common.HexToAddress(settings.IssuanceContractAddress)
		srvc := services.NewTokenTransferService(&settings, producer, usersClient, devicesClient, addr, pdb.DBS)
		err = srvc.TransferUserTokens(ctx, week)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to transfer tokens")
		}
		// need to add in parsing the response to determine success/ failure
	case "listen":
		kconf := sarama.NewConfig()
		kconf.Version = sarama.V2_8_1_0
		kconf.Producer.Return.Successes = true
		kconf.Producer.Partitioner = kafkautil.NewJVMCompatiblePartitioner

		kclient, err := sarama.NewClient(strings.Split(settings.KafkaBrokers, ","), kconf)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
		}

		log.Printf("brokers=%q", settings.KafkaBrokers)
		consumer, err := sarama.NewConsumerGroupFromClient("c1", kclient)
		if err != nil {
			log.Fatal(err)
		}
		defer consumer.Close()
		var wg sync.WaitGroup
		wg.Add(9)
		go services.Consume(consumer, &wg, "c1")
		wg.Wait()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
	default:
		logger.Fatal().Msgf("Unrecognized sub-command %s.", subCommand)
	}
}

func startGRPCServer(settings *config.Settings, dbs func() *database.DBReaderWriter, logger *zerolog.Logger) {
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
