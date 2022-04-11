package services

import (
	"context"
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

// GeetWeekNum calculates the number of the week in which the given time lies for DIMO point
// issuance, which at the time of writing starts at 2022-01-31 05:00 UTC. Indexing is
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

type integrationPointsCalculator struct {
	autoPiID, teslaID, smartcarID string
}

func (i *integrationPointsCalculator) Calculate(integrationIDs []string) int {
	if ContainsString(integrationIDs, i.autoPiID) {
		if ContainsString(integrationIDs, i.smartcarID) {
			return 7000
		} else {
			return 6000
		}
	} else if ContainsString(integrationIDs, i.teslaID) {
		return 4000
	} else if ContainsString(integrationIDs, i.smartcarID) {
		return 1000
	}
	return 0
}

func (t *RewardsTask) createIntegrationPointsCalculator(resp *pb.ListIntegrationsResponse) *integrationPointsCalculator {
	calc := new(integrationPointsCalculator)

	for _, integration := range resp.Integrations {
		switch integration.Vendor {
		case "AutoPi":
			calc.autoPiID = integration.Id
		case "Tesla":
			calc.teslaID = integration.Id
		case "SmartCar":
			calc.smartcarID = integration.Id
		default:
			t.Logger.Warn().Msgf("Unrecognized integration %s with vendor %s", integration.Id, integration.Vendor)
		}
	}

	return calc
}

func (t *RewardsTask) Calculate(issuanceWeek int) error {
	ctx := context.Background()

	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)

	t.Logger.Info().Msgf("Running job for issuance week %d, running from %s to %s", issuanceWeek, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339))

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

	integCalc := t.createIntegrationPointsCalculator(integs)
	vendorToIntegration := make(map[string]string)
	for _, i := range integs.Integrations {
		vendorToIntegration[i.Vendor] = i.Id
	}

	deviceClient := pb.NewUserDeviceServiceClient(conn)

	lastWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek-1)).All(ctx, t.DB().Reader)
	if err != nil {
		return err
	}

	lastWeekByDevice := make(map[string]*models.Reward)
	for _, reward := range lastWeekRewards {
		lastWeekByDevice[reward.UserDeviceID] = reward
	}

	for _, device := range devices {
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

		thisWeek := &models.Reward{
			UserDeviceID:   device.ID,
			IssuanceWeekID: issuanceWeek,
			UserID:         ud.UserId,
		}

		// Streak rewards.
		streakInput := StreakInput{
			ConnectedThisWeek:           true,
			ExistingConnectionStreak:    0,
			ExistingDisconnectionStreak: 0,
		}
		if lastWeek, ok := lastWeekByDevice[device.ID]; ok {
			if lastWeek.UserID != ud.UserId {
				t.Logger.Warn().Str("userDeviceId", ud.Id).Msgf("Device changed ownership from %s to %s, resetting streaks.", lastWeek.UserID, ud.UserId)
			} else {
				streakInput.ExistingConnectionStreak = lastWeek.ConnectionStreak
				streakInput.ExistingDisconnectionStreak = lastWeek.DisconnectionStreak
			}
			delete(lastWeekByDevice, device.ID)
		}

		setStreakFields(thisWeek, ComputeStreak(streakInput))

		// Integration or "connected method" rewards.
		thisWeek.IntegrationIds = device.Integrations
		thisWeek.IntegrationPoints = integCalc.Calculate(device.Integrations)

		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	// We didn't see any data for these remaining devices this week.
	for _, lastWeek := range lastWeekByDevice {
		thisWeek := &models.Reward{
			IssuanceWeekID: issuanceWeek,
			UserDeviceID:   lastWeek.UserDeviceID,
			UserID:         lastWeek.UserID,
		}
		streakInput := StreakInput{
			ConnectedThisWeek:           false,
			ExistingConnectionStreak:    lastWeek.ConnectionStreak,
			ExistingDisconnectionStreak: lastWeek.DisconnectionStreak,
		}
		setStreakFields(thisWeek, ComputeStreak(streakInput))
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

func setStreakFields(reward *models.Reward, streakOutput StreakOutput) {
	reward.ConnectionStreak = streakOutput.ConnectionStreak
	reward.DisconnectionStreak = streakOutput.DisconnectionStreak
	reward.StreakPoints = streakOutput.Points
}
