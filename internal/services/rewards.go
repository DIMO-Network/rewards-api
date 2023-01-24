package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/rewards-api/models"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/volatiletech/null/v8"

	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var baseWeeklyTokens = new(big.Int).Mul(big.NewInt(1_105_000), ether)

var startTime = time.Date(2022, time.January, 31, 5, 0, 0, 0, time.UTC)

var weekDuration = 7 * 24 * time.Hour

// GetWeekNum calculates the number of the week in which the given time lies for DIMO point
// issuance, which at the time of writing starts at 2022-01-31 05:00 UTC. Indexing is
// zero-based.
func GetWeekNum(t time.Time) int {
	sinceStart := t.Sub(startTime)
	weekNum := int(sinceStart.Truncate(weekDuration) / weekDuration)
	return weekNum
}

// GetWeekNumForCron calculates the week number for the current run of the cron job. We expect
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
	Settings        *config.Settings
	Logger          *zerolog.Logger
	DataService     DeviceDataClient
	DB              func() *database.DBReaderWriter
	TransferService Transfer
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
	AutoPiID, TeslaID, SmartcarID string
}

func (i *integrationPointsCalculator) Calculate(integrationIDs []string) int {
	if ContainsString(integrationIDs, i.AutoPiID) {
		if ContainsString(integrationIDs, i.SmartcarID) {
			return 7000
		}
		return 6000
	} else if ContainsString(integrationIDs, i.TeslaID) {
		return 4000
	} else if ContainsString(integrationIDs, i.SmartcarID) {
		return 1000
	}
	return 0
}

func (t *RewardsTask) createIntegrationPointsCalculator(resp *pb_defs.GetIntegrationResponse) *integrationPointsCalculator {
	calc := new(integrationPointsCalculator)

	for _, integration := range resp.Integrations {
		switch integration.Vendor {
		case "AutoPi":
			calc.AutoPiID = integration.Id
		case "Tesla":
			calc.TeslaID = integration.Id
		case "SmartCar":
			calc.SmartcarID = integration.Id
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

	// There shouldn't be anything there. This used to be used when we'd do historical overrides.
	delCount, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).DeleteAll(ctx, t.DB().Writer)
	if err != nil {
		return err
	}

	if delCount != 0 {
		t.Logger.Warn().Int("issuanceWeek", issuanceWeek).Int64("deleted", delCount).Msg("Deleted some existing rows.")
	}

	week := models.IssuanceWeek{
		ID:        issuanceWeek,
		JobStatus: models.IssuanceWeeksJobStatusStarted,
		StartsAt:  weekStart,
		EndsAt:    weekEnd,
	}

	if err := week.Upsert(ctx, t.DB().Writer.DB, true, []string{models.IssuanceWeekColumns.ID}, boil.Whitelist(models.IssuanceWeekColumns.JobStatus), boil.Infer()); err != nil {
		return err
	}

	overrides, err := models.Overrides(models.OverrideWhere.IssuanceWeekID.EQ(issuanceWeek)).All(ctx, t.DB().Reader)
	if err != nil {
		return err
	}

	deviceToOverride := make(map[string]int)
	for _, ov := range overrides {
		deviceToOverride[ov.UserDeviceID] = ov.ConnectionStreak
	}

	// These devices have each sent some signal during the issuance week.
	devices, err := t.DataService.DescribeActiveDevices(weekStart, weekEnd)
	if err != nil {
		return err
	}

	devicesConn, err := grpc.Dial(t.Settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer devicesConn.Close()

	definitionsConn, err := grpc.Dial(t.Settings.DefinitionsAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer definitionsConn.Close()

	definitionsClient := pb_defs.NewDeviceDefinitionServiceClient(definitionsConn)
	integs, err := definitionsClient.GetIntegrations(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	integCalc := t.createIntegrationPointsCalculator(integs)
	vendorToIntegration := make(map[string]string)
	for _, i := range integs.Integrations {
		vendorToIntegration[i.Vendor] = i.Id
	}

	deviceClient := pb_devices.NewUserDeviceServiceClient(devicesConn)

	lastWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek-1)).All(ctx, t.DB().Reader)
	if err != nil {
		return err
	}

	lastWeekByDevice := make(map[string]*models.Reward)
	for _, reward := range lastWeekRewards {
		lastWeekByDevice[reward.UserDeviceID] = reward
	}

	for _, device := range devices {
		ud, err := deviceClient.GetUserDevice(ctx, &pb_devices.GetUserDeviceRequest{Id: device.ID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				t.Logger.Info().Str("userDeviceId", device.ID).Msg("Device was active during the week but was later deleted.")
				continue
			}
			return err
		}

		thisWeek := &models.Reward{
			UserDeviceID:   device.ID,
			IssuanceWeekID: issuanceWeek,
			UserID:         ud.UserId,
		}

		if ud.TokenId == nil {
			t.Logger.Info().Str("userDeviceId", ud.Id).Str("userId", ud.UserId).Msg("Device not minted.")
			continue
		}

		if ud.OptedInAt == nil {
			t.Logger.Info().Str("userDeviceId", ud.Id).Str("userId", ud.UserId).Msg("User has not opted in for this device.")
			continue
		}

		if len(ud.OwnerAddress) != 20 {
			t.Logger.Error().Str("userId", ud.UserId).Bytes("address", ud.OwnerAddress).Msg("User has minted a car but has no owner address?")
			continue
		}

		if ContainsString(device.Integrations, integCalc.AutoPiID) {
			if ud.AftermarketDeviceTokenId == nil {
				t.Logger.Info().Str("userDeviceId", ud.Id).Msg("AutoPi activity but not paired on-chain.")

				filtered := []string{}

				for _, integ := range device.Integrations {
					if integ != integCalc.AutoPiID {
						filtered = append(filtered, integ)
					}
				}

				// If the only active connection was an AutoPi, and this was not paired on chain,
				// it doesn't count.
				if len(filtered) == 0 {
					continue
				}

				device.Integrations = filtered
			} else {
				thisWeek.AftermarketTokenID = types.NewNullDecimal(new(decimal.Big).SetUint64(*ud.AftermarketDeviceTokenId))
			}
		}

		thisWeek.UserDeviceTokenID = types.NewNullDecimal(new(decimal.Big).SetUint64(*ud.TokenId))
		thisWeek.UserEthereumAddress = null.StringFrom(common.BytesToAddress(ud.OwnerAddress).Hex())

		var streak StreakOutput

		if connStreak, ok := deviceToOverride[device.ID]; ok {
			streak = FakeStreak(connStreak)
			delete(deviceToOverride, device.ID)
		} else {
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

			streak = ComputeStreak(streakInput)
		}

		setStreakFields(thisWeek, streak)

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
		streak := ComputeStreak(streakInput)
		setStreakFields(thisWeek, streak)
		if err := thisWeek.Insert(ctx, t.DB().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	if len(deviceToOverride) != 0 {
		t.Logger.Warn().Interface("overrides", deviceToOverride).Msg("Unused overrides.")
	}

	st := storage.NewDB(t.DB)
	err = st.AssignTokens(ctx, issuanceWeek, baseWeeklyTokens)
	if err != nil {
		return fmt.Errorf("failed to convert points to tokens: %w", err)
	}

	err = t.TransferService.TransferUserTokens(ctx, issuanceWeek)
	if err != nil {
		return fmt.Errorf("failed to submit transfers: %w", err)
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
