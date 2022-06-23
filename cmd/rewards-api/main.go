package main

import (
	"context"
	"errors"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"os"
	"strconv"
	"time"

	_ "github.com/DIMO-Network/rewards-api/docs"
	"github.com/DIMO-Network/rewards-api/internal/controllers"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/DIMO-Network/shared"
	pb "github.com/DIMO-Network/shared/api/devices"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
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

		conn, err := grpc.Dial(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to create devices API client.")
		}
		defer conn.Close()

		integClient := pb.NewIntegrationServiceClient(conn)
		deviceClient := pb.NewUserDeviceServiceClient(conn)

		dataClient := services.NewDeviceDataClient(&settings)

		rewardsController := controllers.RewardsController{
			DB:            pdb.DBS,
			Logger:        &logger,
			IntegClient:   integClient,
			DevicesClient: deviceClient,
			DataClient:    dataClient,
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
			logger.Fatal().Err(err).Msg("Failed to calculate rewards for week %d.")
		}
	default:
		logger.Fatal().Msgf("Unrecognized sub-command %s.", subCommand)
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError // Default.

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": err.Error(),
	})
}
