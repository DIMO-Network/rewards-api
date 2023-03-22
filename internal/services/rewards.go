package services

import (
	"context"
	"fmt"
	"math/big"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared"
	pb_devices "github.com/DIMO-Network/shared/api/devices"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/volatiletech/null/v8"
	"golang.org/x/exp/slices"

	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var baseWeeklyTokens = new(big.Int).Mul(big.NewInt(1_105_000), ether)

var startTime = time.Date(2022, time.January, 31, 5, 0, 0, 0, time.UTC)

var weekDuration = 7 * 24 * time.Hour

type BaselineClient struct {
	TransferService *TransferService
	DataService     DeviceActivityClient
	DevicesClient   DevicesClient
	DefsClient      IntegrationsGetter
	ContractAddress common.Address
	Week            int
	Logger          *zerolog.Logger
}

type DeviceActivityClient interface {
	DescribeActiveDevices(start, end time.Time) ([]*DeviceData, error)
}

type IntegrationsGetter interface {
	GetIntegrations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb_defs.GetIntegrationResponse, error)
}

type DevicesClient interface {
	GetUserDevice(ctx context.Context, in *pb_devices.GetUserDeviceRequest, opts ...grpc.CallOption) (*pb_devices.UserDevice, error)
}

type integrationPointsCalculator struct {
	AutoPiID   string
	TeslaID    string
	SmartcarID string
}

func NewBaselineRewardService(
	settings *config.Settings,
	transferService *TransferService,
	dataService DeviceActivityClient,
	devicesClient DevicesClient,
	defsClient IntegrationsGetter,
	week int,
	logger *zerolog.Logger,

) *BaselineClient {

	return &BaselineClient{
		TransferService: transferService,
		DataService:     dataService,
		DevicesClient:   devicesClient,
		DefsClient:      defsClient,
		ContractAddress: common.HexToAddress(settings.IssuanceContractAddress),
		Week:            week,
		Logger:          logger,
	}
}

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

func (i *integrationPointsCalculator) Calculate(integrationIDs []string) int {
	// Only blessed combination.
	if slices.Contains(integrationIDs, i.AutoPiID) {
		if slices.Contains(integrationIDs, i.SmartcarID) {
			return 7000
		}
		return 6000
	} else if slices.Contains(integrationIDs, i.TeslaID) {
		return 4000
	} else if slices.Contains(integrationIDs, i.SmartcarID) {
		return 1000
	}
	return 0
}

func (t *BaselineClient) createIntegrationPointsCalculator(resp *pb_defs.GetIntegrationResponse) *integrationPointsCalculator {
	var calc integrationPointsCalculator

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

	return &calc
}

func (t *BaselineClient) Calculate(issuanceWeek int) error {
	ctx := context.Background()

	weekStart := NumToWeekStart(issuanceWeek)
	weekEnd := NumToWeekEnd(issuanceWeek)

	t.Logger.Info().Msgf("Running job for issuance week %d, running from %s to %s", issuanceWeek, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339))

	// There shouldn't be anything there. This used to be used when we'd do historical overrides.
	delCount, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).DeleteAll(ctx, t.TransferService.db.DBS().Writer)
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

	if err := week.Upsert(ctx, t.TransferService.db.DBS().Writer, true, []string{models.IssuanceWeekColumns.ID}, boil.Whitelist(models.IssuanceWeekColumns.JobStatus), boil.Infer()); err != nil {
		return err
	}

	overrides, err := models.Overrides(models.OverrideWhere.IssuanceWeekID.EQ(issuanceWeek)).All(ctx, t.TransferService.db.DBS().Reader)
	if err != nil {
		return err
	}

	deviceToOverride := map[string]int{}
	for _, ov := range overrides {
		deviceToOverride[ov.UserDeviceID] = ov.ConnectionStreak
	}

	// These describe the active integrations for each device active this week.
	deviceActivityRecords, err := t.DataService.DescribeActiveDevices(weekStart, weekEnd)
	if err != nil {
		return err
	}

	integs, err := t.DefsClient.GetIntegrations(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	integCalc := t.createIntegrationPointsCalculator(integs)

	lastWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek-1)).All(ctx, t.TransferService.db.DBS().Reader)
	if err != nil {
		return err
	}

	lastWeekByDevice := make(map[string]*models.Reward)
	for _, reward := range lastWeekRewards {
		lastWeekByDevice[reward.UserDeviceID] = reward
	}

	seenVINs, err := models.Vins().All(ctx, t.TransferService.db.DBS().Reader)
	if err != nil {
		return err
	}

	seenVINSet := shared.NewStringSet()
	for _, vinRec := range seenVINs {
		seenVINSet.Add(vinRec.Vin)
	}

	for _, deviceActivity := range deviceActivityRecords {
		logger := t.Logger.With().Str("userDeviceId", deviceActivity.ID).Logger()

		ud, err := t.DevicesClient.GetUserDevice(ctx, &pb_devices.GetUserDeviceRequest{Id: deviceActivity.ID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				logger.Info().Msg("Device was active during the week but was later deleted.")
				continue
			}
			return err
		}

		logger = logger.With().Str("userId", ud.UserId).Logger()

		if ud.TokenId == nil {
			logger.Info().Msg("Device not minted.")
			continue
		}

		if ud.OptedInAt == nil {
			logger.Info().Msg("User has not opted in for this device.")
			continue
		}

		if len(ud.OwnerAddress) != 20 {
			logger.Info().Msg("User has minted a car but has no owner address?")
			continue
		}

		thisWeek := &models.Reward{
			UserDeviceID:        deviceActivity.ID,
			IssuanceWeekID:      issuanceWeek,
			UserID:              ud.UserId,
			UserDeviceTokenID:   types.NewNullDecimal(new(decimal.Big).SetUint64(*ud.TokenId)),
			UserEthereumAddress: null.StringFrom(common.BytesToAddress(ud.OwnerAddress).Hex()),
		}

		validIntegrations := deviceActivity.Integrations // Guaranteed to be non-empty at this point.

		if ind := slices.Index(validIntegrations, integCalc.AutoPiID); ind != -1 {
			if ud.AftermarketDeviceTokenId == nil {
				if len(validIntegrations) == 1 {
					logger.Info().Msg("AutoPi connected but not paired-onchain; no other active integrations.")
					continue
				} else {
					validIntegrations = slices.Delete(validIntegrations, ind, ind+1)
					logger.Info().Msg("AutoPi connected but not paired on-chain; there are other integrations.")
				}
			} else {
				thisWeek.AftermarketTokenID = types.NewNullDecimal(new(decimal.Big).SetUint64(*ud.AftermarketDeviceTokenId))
			}
		}

		// At this point we are certain that the owner should receive tokens.
		thisWeek.IntegrationIds = validIntegrations
		thisWeek.IntegrationPoints = integCalc.Calculate(validIntegrations)

		var streak StreakOutput

		if connStreak, ok := deviceToOverride[deviceActivity.ID]; ok {
			logger.Info().Int("connectionStreak", connStreak).Msg("Override for active device.")
			streak = FakeStreak(connStreak)
			delete(deviceToOverride, deviceActivity.ID)
		} else {
			// Streak rewards.
			streakInput := StreakInput{
				ConnectedThisWeek:           true,
				ExistingConnectionStreak:    0,
				ExistingDisconnectionStreak: 0,
			}
			if lastWeek, ok := lastWeekByDevice[deviceActivity.ID]; ok {
				streakInput.ExistingConnectionStreak = lastWeek.ConnectionStreak
				streakInput.ExistingDisconnectionStreak = lastWeek.DisconnectionStreak

			}

			streak = ComputeStreak(streakInput)
		}
		setStreakFields(thisWeek, streak)

		// Anything left in this map is considered disconnected.
		// This is a no-op if the device doesn't have a record from last week.
		delete(lastWeekByDevice, deviceActivity.ID)

		// If this VIN has never earned before, make note of that.
		// Used by referrals, not this job.
		if ud.Vin != nil && len(*ud.Vin) == 17 && !seenVINSet.Contains(*ud.Vin) {
			seenVINSet.Add(*ud.Vin)
			vinRec := models.Vin{
				Vin:              *ud.Vin,
				FirstWeekEarning: issuanceWeek,
			}
			if err := vinRec.Insert(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
				return err
			}
			thisWeek.NewVin = true
		}

		if err := thisWeek.Insert(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
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
		if err := thisWeek.Insert(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	if len(deviceToOverride) != 0 {
		t.Logger.Warn().Interface("overrides", deviceToOverride).Msg("Unused overrides.")
	}

	st := storage.DBStorage{DBS: t.TransferService.db}
	err = st.AssignTokens(ctx, issuanceWeek, baseWeeklyTokens)
	if err != nil {
		return fmt.Errorf("failed to convert points to tokens: %w", err)
	}

	return nil
}

func (t *BaselineClient) BaselineIssuance(issuanceWeek int) error {
	ctx := context.Background()

	err := t.Calculate(t.Week)
	if err != nil {
		return fmt.Errorf("failed to calculate and assign baseline tokens: %w", err)
	}

	err = t.transfer(ctx)
	if err != nil {
		return fmt.Errorf("failed to submit baseline token transfers: %w", err)
	}

	week, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(issuanceWeek)).One(ctx, t.TransferService.db.DBS().Writer)
	if err != nil {
		return err
	}

	week.JobStatus = models.IssuanceWeeksJobStatusFinished
	if _, err := week.Update(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func setStreakFields(reward *models.Reward, streakOutput StreakOutput) {
	reward.ConnectionStreak = streakOutput.ConnectionStreak
	reward.DisconnectionStreak = streakOutput.DisconnectionStreak
	reward.StreakPoints = streakOutput.Points
}
