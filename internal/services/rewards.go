package services

import (
	"context"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var startTime = time.Date(2022, time.January, 31, 5, 0, 0, 0, time.UTC)

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

func NumToWeekStart(n int) time.Time {
	return startTime.Add(time.Duration(n) * weekDuration)
}

func NumToWeekEnd(n int) time.Time {
	return startTime.Add(time.Duration(n+1) * weekDuration)
}

type RewardsTask struct {
	Settings    *config.Settings
	Logger      *zerolog.Logger
	DataService DeviceDataClient
	DB          func() *database.DBReaderWriter
}

type ConnectionMethod struct {
	DevicesAPIVendor string
	DBConstant       string
	Points           int
}

func ContainsString(v []string, x string) bool {
	for _, y := range v {
		if y == x {
			return true
		}
	}
	return false
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

	// These devices have each sent some signal during the issuance week.
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

	vendorToIntegration := make(map[string]string)
	for _, i := range integs.Integrations {
		vendorToIntegration[i.Vendor] = i.Id
	}

	smartcarIntegrationID, ok := vendorToIntegration["SmartCar"]
	if !ok {
		t.Logger.Warn().Msg("Device service did not return a Smartcar integration.")
	}
	teslaIntegrationID, ok := vendorToIntegration["Tesla"]
	if !ok {
		t.Logger.Warn().Msg("Device service did not return a Tesla integration.")
	}
	autoPiIntegrationID := vendorToIntegration["AutoPi"]
	if !ok {
		return errors.New("Devices service did not return an AutoPi integration.")
	}

	fmt.Println(smartcarIntegrationID, teslaIntegrationID, autoPiIntegrationID)

	lastWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek-1)).All(ctx, t.DB().Reader)
	if err != nil {
		return err
	}

	lastWeekByDevice := make(map[string]*models.Reward)
	for _, reward := range lastWeekRewards {
		lastWeekByDevice[reward.UserDeviceID] = reward
	}

	total := len(devices)
	for i, device := range devices {
		fmt.Printf("%d/%d\n", i+1, total)

		ud, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{Id: device.ID})
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.NotFound {
					t.Logger.Info().Str("userDeviceId", device.ID).Msg("Device was active during the week but was later deleted.")
					continue
				} else {
					return err
				}
			} else {
				return err
			}
		}

		thisWeek := models.Reward{
			UserDeviceID:   device.ID,
			IssuanceWeekID: issuanceWeek,
			UserID:         ud.UserId,
		}

		// Streak rewards.
		x := Input{ConnectedThisWeek: true}
		lastWeek, ok := lastWeekByDevice[device.ID]
		if ok {
			x.ExistingDisconnectionStreak = lastWeek.DisconnectionStreak
			x.ExistingConnectionStreak = lastWeek.EffectiveConnectionStreak
			delete(lastWeekByDevice, device.ID)
		}

		o := ComputeStreak(x)
		thisWeek.EffectiveConnectionStreak = o.ConnectionStreak
		thisWeek.DisconnectionStreak = o.DisconnectionStreak
		thisWeek.StreakPoints = o.Points

		// Integration or "connected method" rewards.
		thisWeek.IntegrationIds = device.Integrations

		if ContainsString(device.Integrations, autoPiIntegrationID) {
			if ContainsString(device.Integrations, smartcarIntegrationID) {
				thisWeek.IntegrationPoints = 7000
			} else {
				thisWeek.IntegrationPoints = 6000
			}
		} else if ContainsString(device.Integrations, teslaIntegrationID) {
			thisWeek.IntegrationPoints = 4000
		} else if ContainsString(device.Integrations, smartcarIntegrationID) {
			thisWeek.IntegrationPoints = 1000
		}

		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	for _, lastWeek := range lastWeekByDevice {
		x := Input{
			ExistingDisconnectionStreak: lastWeek.DisconnectionStreak,
			ExistingConnectionStreak:    lastWeek.EffectiveConnectionStreak,
			ConnectedThisWeek:           false,
		}
		o := ComputeStreak(x)
		thisWeek := models.Reward{
			UserDeviceID:              lastWeek.UserDeviceID,
			IssuanceWeekID:            issuanceWeek,
			UserID:                    lastWeek.UserID,
			EffectiveConnectionStreak: o.ConnectionStreak,
			DisconnectionStreak:       o.DisconnectionStreak,
		}
		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	week.JobStatus = models.IssuanceWeeksJobStatusFinished
	if _, err := week.Update(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}

	return nil
}
