package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb "github.com/DIMO-Network/shared/api/devices"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var startTime = time.Date(2022, time.January, 31, 5, 0, 0, 0, time.UTC)

var weekDuration = 7 * 24 * time.Hour

type AllPointsDistributed struct {
	DistributedPoints *big.Int `boil:"distributed_points"`
}

type DeviceTotalPoints struct {
	DevicePoints int64 `boil:"device_points"`
}

// GeetWeekNum calculates the number of the week in which the given time lies for DIMO point
// issuance, which at the time of writing starts at 2022-01-31 05:00 UTC. Indexing is
// zero-based.
func GetWeekNum(t time.Time) int {
	sinceStart := t.Sub(startTime)
	weekNum := int(sinceStart.Truncate(weekDuration) / weekDuration)
	return weekNum
}

// GeetWeekNumForCron calculates the week number for the current run of the cron job. We expect
// the job to run every Monday at 05:00 UTC, but due to skew we just round the time.
func GetWeekNumForCron(t time.Time) int {
	sinceStart := t.Sub(startTime)
	weekNum := int(sinceStart.Round(weekDuration) / weekDuration)
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
		}
		return 6000
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
	var totalPointsDistributed int64

	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)

	t.Logger.Info().Msgf("Running job for issuance week %d, running from %s to %s", issuanceWeek, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339))

	if _, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
		models.RewardWhere.Override.EQ(false),
	).DeleteAll(ctx, t.DB().Writer); err != nil {
		return err
	}

	dimo := WeeklyTokenAllocation(issuanceWeek)
	week := models.IssuanceWeek{
		ID:                    issuanceWeek,
		JobStatus:             models.IssuanceWeeksJobStatusStarted,
		StartsAt:              weekStart,
		EndsAt:                weekEnd,
		WeeklyTokenAllocation: types.NewNullDecimal(new(decimal.Big).SetBigMantScale(dimo, 0)),
	}

	if err := week.Upsert(ctx, t.DB().Writer.DB, true, []string{models.IssuanceWeekColumns.ID}, boil.Whitelist(models.IssuanceWeekColumns.JobStatus), boil.Infer()); err != nil {
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

	overrides, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
		models.RewardWhere.Override.EQ(true),
	).All(ctx, t.DB().Reader.DB)
	if err != nil {
		return err
	}

	devicesOverriddenThisWeek := shared.NewStringSet()

	for _, override := range overrides {
		if len(override.IntegrationIds) == 0 {
			t.Logger.Warn().Str("userDeviceId", override.UserDeviceID).Msg("Override had no integrations.")
			continue
		}
		devicesOverriddenThisWeek.Add(override.UserDeviceID)

		ud, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{Id: override.UserDeviceID})
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.NotFound {
					t.Logger.Info().Str("userDeviceId", override.UserDeviceID).Msg("Override present for deleted device.")
					continue
				} else {
					return err
				}
			} else {
				return err
			}
		}

		override.UserID = ud.UserId

		streakInput := StreakInput{
			ConnectedThisWeek:           true,
			ExistingConnectionStreak:    0,
			ExistingDisconnectionStreak: 0,
		}

		if lastWeek, ok := lastWeekByDevice[override.UserDeviceID]; ok {
			if lastWeek.UserID != ud.UserId {
				t.Logger.Warn().Str("userDeviceId", ud.Id).Msgf("Device changed ownership from %s to %s, resetting streaks.", lastWeek.UserID, ud.UserId)
			} else {
				streakInput.ExistingConnectionStreak = lastWeek.ConnectionStreak
				streakInput.ExistingDisconnectionStreak = lastWeek.DisconnectionStreak
			}
			delete(lastWeekByDevice, override.UserDeviceID)
		}
		streak := ComputeStreak(streakInput)
		setStreakFields(override, streak)
		totalPointsDistributed += int64(streak.Points)

		override.IntegrationPoints = integCalc.Calculate(override.IntegrationIds)
		totalPointsDistributed += int64(override.IntegrationPoints)

		if _, err := override.Update(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	for _, device := range devices {
		if devicesOverriddenThisWeek.Contains(device.ID) {
			continue
		}

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
		streak := ComputeStreak(streakInput)
		setStreakFields(thisWeek, streak)
		totalPointsDistributed += int64(streak.Points)

		// Integration or "connected method" rewards.
		thisWeek.IntegrationIds = device.Integrations
		thisWeek.IntegrationPoints = integCalc.Calculate(device.Integrations)
		totalPointsDistributed += int64(thisWeek.IntegrationPoints)

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
		streak := ComputeStreak(streakInput)
		setStreakFields(thisWeek, streak)
		totalPointsDistributed += int64(streak.Points)

		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	week.JobStatus = models.IssuanceWeeksJobStatusPointsAllocated
	week.PointsDistributed = null.Int64From(totalPointsDistributed)
	if _, err := week.Update(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func (t *RewardsTask) Allocate(issuanceWeek int) error {
	ctx := context.Background()
	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)
	t.Logger.Info().Msgf("Running token allocation for issuance week %d, running from %s to %s", issuanceWeek, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339))
	deviceRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).All(ctx, t.DB().Reader)
	if err != nil {
		return err
	}
	distribution, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(issuanceWeek)).One(ctx, t.DB().Reader)
	if err != nil {
		return err
	}
	tknBytes, err := distribution.WeeklyTokenAllocation.MarshalText()
	if err != nil {
		return err
	}
	distributedTokens := new(big.Int)
	distributedTokens, ok := distributedTokens.SetString(string(tknBytes), 10)
	if !ok {
		fmt.Println("SetString: error")
		return nil
	}
	distribution.JobStatus = models.IssuanceWeeksJobStatusBeginTokenDistribution
	if _, err := distribution.Update(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}
	for _, device := range deviceRewards {
		devicePoints := device.IntegrationPoints + device.StreakPoints
		deviceTokens := CalculateTokenAllocation(devicePoints, distribution.PointsDistributed.Int64, distributedTokens)

		update := models.TokenAllocation{
			IssuanceWeekID: issuanceWeek,
			UserDeviceID:   device.UserDeviceID,
			Tokens:         types.NewNullDecimal(new(decimal.Big).SetBigMantScale(deviceTokens, 0)),
			WeekStart:      weekStart,
			WeekEnd:        weekEnd,
		}

		if err := update.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}
	distribution.JobStatus = models.IssuanceWeeksJobStatusFinished
	if _, err := distribution.Update(ctx, t.DB().Writer, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func setStreakFields(reward *models.Reward, streakOutput StreakOutput) {
	reward.ConnectionStreak = streakOutput.ConnectionStreak
	reward.DisconnectionStreak = streakOutput.DisconnectionStreak
	reward.StreakPoints = streakOutput.Points
}
