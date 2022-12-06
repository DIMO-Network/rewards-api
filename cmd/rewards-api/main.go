package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	_ "github.com/DIMO-Network/rewards-api/docs"
	"github.com/DIMO-Network/rewards-api/internal/api"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/controllers"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/shared"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	pb_rewards "github.com/DIMO-Network/shared/api/rewards"
	"github.com/Shopify/sarama"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/burdiyan/kafkautil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title                       DIMO Rewards API
// @version                     1.0
// @BasePath                    /v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("app", "rewards-api").Logger()

	ctx, cancel := context.WithCancel(context.Background())
	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load settings.")
	}

	if len(os.Args) == 1 {
		monApp := serveMonitoring(settings.MonitoringPort, &logger)
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
		defer definitionsConn.Close()

		definitionsClient := pb_defs.NewDeviceDefinitionServiceClient(definitionsConn)
		deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)

		dataClient := services.NewDeviceDataClient(&settings)

		rewardsController := controllers.RewardsController{
			DB:                pdb.DBS,
			Logger:            &logger,
			DefinitionsClient: definitionsClient,
			DevicesClient:     deviceClient,
			DataClient:        dataClient,
			Settings:          &settings,
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

		go startGRPCServer(&settings, pdb.DBS, &logger)

		// Start metatransaction listener
		if settings.Environment != "prod" {
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
			statusProc, err := services.NewStatusProcessor(pdb.DBS, &logger)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to create transaction status processor.")
			}

			go func() {
				err := services.Consume(ctx, consumer, &settings, statusProc)
				if err != nil {
					logger.Fatal().Err(err).Send()
				}
			}()
		}

		logger.Info().Msgf("Starting HTTP server on port %s.", settings.Port)
		if err := app.Listen(":" + settings.Port); err != nil {
			logger.Fatal().Err(err).Msgf("Fiber server failed.")
		}

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, os.Interrupt)

		sig := <-sigterm
		logger.Info().Str("signal", sig.String()).Msg("Received signal, terminating.")

		cancel()
		monApp.Shutdown()

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

		var transferService services.Transfer

		if settings.Environment != "prod" {
			kclient, err := createKafkaClient(&settings)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to create Kafka client.")
			}

			producer, err := sarama.NewSyncProducerFromClient(kclient)
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to create Kafka producer.")
			}

			addr := common.HexToAddress(settings.IssuanceContractAddress)
			transferService = services.NewTokenTransferService(&settings, producer, addr, pdb.DBS)
		}

		task := services.RewardsTask{
			Settings:        &settings,
			DataService:     services.NewDeviceDataClient(&settings),
			DB:              pdb.DBS,
			Logger:          &logger,
			TransferService: transferService,
		}
		if err := task.Calculate(week); err != nil {
			logger.Fatal().Err(err).Int("issuanceWeek", week).Msg("Failed to calculate and/or transfer rewards.")
		}
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

func serveMonitoring(port string, logger *zerolog.Logger) *fiber.App {
	monApp := fiber.New(fiber.Config{DisableStartupMessage: true})

	// Health check.
	monApp.Get("/", func(c *fiber.Ctx) error { return nil })
	monApp.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	go func() {
		if err := monApp.Listen(":" + port); err != nil {
			logger.Fatal().Err(err).Str("port", port).Msg("Failed to start monitoring web server.")
		}
	}()

	logger.Info().Str("port", port).Msg("Started monitoring web server.")

	return monApp
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
