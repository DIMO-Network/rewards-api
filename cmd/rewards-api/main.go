package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/services"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	settings, err := config.LoadConfig("settings.yaml")
	if err != nil {
		os.Exit(1)
	}

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "rewards-api").
		Logger()

	if len(os.Args) == 1 {
		logger.Fatal().Msg("Sub-command required.")
	}

	switch subCommand := os.Args[1]; subCommand {
	case "migrate":
		migrateDatabase(logger, settings)
	case "getweek":
		if len(os.Args) < 3 {
			logger.Fatal().Msg("Date string required.")
		}
		dateString := os.Args[2]
		t, err := time.Parse(time.RFC3339, dateString)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not parse date string %v.", dateString)
		}
		fmt.Printf("Issuance week: %d\n", services.GetWeekNum(t))
	case "calc":
		if len(os.Args) < 3 {
			logger.Fatal().Msg("Issuance week required.")
		}
		weekStr := os.Args[2]
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			logger.Fatal().Err(err).Msg("Could not parse week number.")
		}
		if week < 0 {
			logger.Fatal().Msgf("Negative week number %d.", week)
		}
		pdb := database.NewDbConnectionFromSettings(ctx, settings)
		totalTime := 0
		for !pdb.IsReady() {
			if totalTime > 30 {
				logger.Fatal().Msg("could not connect to postgres after 30 seconds")
			}
			time.Sleep(time.Second)
			totalTime++
		}
		task := services.RewardsTask{
			Settings:    settings,
			DataService: services.NewDeviceDataClient(settings),
			DB:          pdb.DBS,
		}
		if err := task.Calculate(week); err != nil {
			logger.Fatal().Err(err).Msg("Failed to calculate rewards for week %d.")
		}
	default:
		logger.Fatal().Msgf("Unrecognized sub-command %s.", subCommand)
	}
}
