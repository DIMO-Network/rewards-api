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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var startTime = time.Date(2022, time.February, 28, 5, 0, 0, 0, time.UTC)

var weekDuration = 7 * 24 * time.Hour

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
	DataService DeviceDataClient
	DB          func() *database.DBReaderWriter
}

func (t *RewardsTask) Calculate(issuanceWeek int) error {
	ctx := context.Background()

	// WARNING
	// Not production ready, obviously.
	if _, err := models.DeviceWeekRewards(models.DeviceWeekRewardWhere.IssuanceWeekID.EQ(issuanceWeek)).DeleteAll(ctx, t.DB().Writer); err != nil {
		return err
	}

	if _, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(issuanceWeek)).DeleteAll(ctx, t.DB().Writer); err != nil {
		return err
	}

	week := models.IssuanceWeek{
		ID:        issuanceWeek,
		JobStatus: models.IssuanceWeekJobStatusStarted,
	}

	if err := week.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}

	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)

	devIDs, err := t.DataService.ListActiveUserDeviceIDs(weekStart, weekEnd)
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

	integMap := make(map[string]string)
	for _, i := range integs.Integrations {
		integMap[i.Id] = i.Vendor
	}

	evDDCounts := make(map[string]int)

	total := len(devIDs)
	for i, devID := range devIDs {
		fmt.Printf("%d/%d\n", i+1, total)

		ud, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{Id: devID})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if ud.VinIdentifier == nil {
			continue
		}

		thisWeek := models.DeviceWeekReward{
			UserDeviceID:       devID,
			IssuanceWeekID:     issuanceWeek,
			DeviceDefinitionID: ud.DeviceDefinitionId,
			Vin:                *ud.VinIdentifier,
			UserID:             ud.UserId,
		}
		driven, err := t.DataService.GetMilesDriven(devID, weekStart, weekEnd)
		if err != nil {
			return err
		}

		thisWeek.MilesDriven = driven

		var x Input
		last, err := models.FindDeviceWeekReward(ctx, t.DB().Reader, issuanceWeek-1, devID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			x.ExistingDisconnectionStreak = last.WeeksDisconnectedStreak
			x.ExistingConnectionStreak = last.WeeksConnectedStreak
		}

		if driven >= 10 {
			x.ConnectedThisWeek = true
			thisWeek.Connected = true
		}

		o := ComputeStreak(x)
		thisWeek.WeeksConnectedStreak = o.ConnectionStreak
		thisWeek.WeeksDisconnectedStreak = o.DisconnectionStreak
		thisWeek.StreakPoints = o.Points

		if thisWeek.Connected {
			di, err := t.DataService.ListActiveIntegrations(devID, weekStart, weekEnd)
			if err != nil {
				panic(err)
			}

			methodReward := 0
			for _, d := range di {
				vend, ok := integMap[d]
				if !ok {
					fmt.Println("Unknown integration d")
					continue
				}
				switch vend {
				case "Tesla":
					methodReward = 4000
				case "SmartCar":
					methodReward = 1000
				}
			}

			thisWeek.ConnectionMethodPoints = methodReward

			electric, err := t.DataService.UsesElectricity(devID, weekStart, weekEnd)
			if err != nil {
				return err
			}
			if electric {
				gas, err := t.DataService.UsesFuel(devID, weekStart, weekEnd)
				if err != nil {
					return err
				}

				if gas {
					thisWeek.EngineTypePoints = 3000
				} else {
					thisWeek.EngineTypePoints = 5000
				}

				ud, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{Id: devID})
				if err != nil {
					fmt.Println(err)
				} else {
					thisWeek.DeviceDefinitionID = ud.DeviceDefinitionId
					evDDCounts[ud.DeviceDefinitionId] += 1
				}
			}
		}

		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	for dd, ct := range evDDCounts {
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
			models.DeviceWeekRewardWhere.IssuanceWeekID.EQ(issuanceWeek),
			models.DeviceWeekRewardWhere.DeviceDefinitionID.EQ(dd),
			models.DeviceWeekRewardWhere.EngineTypePoints.GT(0),
		).UpdateAll(ctx, t.DB().Writer.DB, models.M{
			"rarity_points": bonus,
		})
		if err != nil {
			panic(err)
		}
	}

	return nil
}
