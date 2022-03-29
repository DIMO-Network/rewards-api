package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	pb "github.com/DIMO-Network/shared/api/devices"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var startTime = time.Date(2022, time.February, 28, 5, 0, 0, 0, time.UTC)

var weekDuration = 7 * 24 * time.Hour

const miToKmFactor = 1.609344

// GeetWeekNum calculates the number of the week in which the given time lies for DIMO point
// issuance, which at the time of writing starts at 05:00 March 14, 2022 UTC. Indexing is
// zero-based.
func GetWeekNum(t time.Time) int {
	sinceStart := t.Sub(startTime)
	weekNum := int(sinceStart.Truncate(weekDuration) / weekDuration)
	return weekNum
}

type RewardsTask struct {
	Settings    *config.Settings
	Logger      *zerolog.Logger
	DataService DeviceDataClient
	DB          func() *database.DBReaderWriter
}

func (t *RewardsTask) Calculate(issuanceWeek int) error {
	ctx := context.Background()

	// Not production ready, obviously.
	if _, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).DeleteAll(ctx, t.DB().Writer); err != nil {
		return err
	}

	if _, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(issuanceWeek)).DeleteAll(ctx, t.DB().Writer); err != nil {
		return err
	}

	week := models.IssuanceWeek{
		ID:        issuanceWeek,
		JobStatus: models.IssuanceWeeksJobStatusStarted,
	}

	if err := week.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}

	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)

	devices, err := t.DataService.DescribeActiveDevices(weekStart, weekEnd)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(t.Settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	integClient := pb.NewIntegrationServiceClient(conn)
	integs, err := integClient.ListIntegrations(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	deviceClient := pb.NewUserDeviceServiceClient(conn)

	integrationVendors := make(map[string]string)
	for _, i := range integs.Integrations {
		integrationVendors[i.Id] = i.Vendor
	}

	greenDeviceCounts := make(map[string]int)

	total := len(devices)
	for i, device := range devices {
		fmt.Printf("%d/%d\n", i+1, total)

		ud, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{Id: device.ID})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if ud.VinIdentifier == nil {
			continue
		}

		thisWeek := models.Reward{
			UserDeviceID:       device.ID,
			IssuanceWeekID:     issuanceWeek,
			DeviceDefinitionID: ud.DeviceDefinitionId,
			Vin:                *ud.VinIdentifier,
			UserID:             ud.UserId,
			DrivenKM:           device.Driven,
		}

		var x Input
		last, err := models.FindReward(ctx, t.DB().Reader, issuanceWeek-1, device.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			x.ExistingDisconnectionStreak = last.DisconnectionStreak
			x.ExistingConnectionStreak = last.ConnectionStreak
		}

		connected := false
		if device.Driven >= 10*miToKmFactor {
			x.ConnectedThisWeek = true
			connected = true
		}

		o := ComputeStreak(x)
		thisWeek.ConnectionStreak = o.ConnectionStreak
		thisWeek.DisconnectionStreak = o.DisconnectionStreak

		if connected {
			device.Integrations
			for _, d := range di {
				vend, ok := integrationVendors[d]
				if !ok {
					t.Logger.Warn().Msgf("Unknown integration %d.", d)
					continue
				}
				// TODO: Allow an AutoPi + SmartCar combination.
				switch vend {
				case "AutoPi":
					methodReward = 6000
				case "Tesla":
					methodReward = 4000
				case "SmartCar":
					methodReward = 1000
				}
			}

			electric, err := t.DataService.UsesElectricity(device, weekStart, weekEnd)
			if err != nil {
				return err
			}
			if electric {
				gas, err := t.DataService.UsesFuel(device, weekStart, weekEnd)
				if err != nil {
					return err
				}

				if gas {
					thisWeek.EngineTypePoints = 3000
				} else {
					thisWeek.EngineTypePoints = 5000
				}

				ud, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{Id: device})
				if err != nil {
					fmt.Println(err)
				} else {
					thisWeek.DeviceDefinitionID = ud.DeviceDefinitionId
					greenDeviceCounts[ud.DeviceDefinitionId] += 1
				}
			}
		}

		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	for dd, ct := range greenDeviceCounts {
		bonus := 0
		switch {
		case ct <= 10:
			bonus = 5000
		case ct <= 50:
			bonus = 3000
		case ct <= 200:
			bonus = 1000
		default:
			continue
		}
		_, err := models.DeviceWeekRewards(
			models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
			models.RewardWhere.DeviceDefinitionID.EQ(dd),
		).UpdateAll(ctx, t.DB().Writer.DB, models.M{
			"rarity_points": bonus,
		})
		if err != nil {
			panic(err)
		}
	}

	week.JobStatus = models.IssuanceWeeksJobStatusFinished
	if _, err := week.Update(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}

	return nil
}
