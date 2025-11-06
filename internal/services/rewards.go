//go:generate mockgen -source=./rewards.go -destination=rewards_mock_test.go -package=services
package services

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/constants"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/services/identity"
	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/DIMO-Network/set"
	pb_tesla "github.com/DIMO-Network/tesla-oracle/pkg/grpc"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaselineClient struct {
	TransferService    *TransferService
	DataService        DeviceActivityClient
	DevicesClient      DevicesClient
	ContractAddress    common.Address
	Week               int
	Logger             *zerolog.Logger
	FirstAutomatedWeek int
	IdentityClient     IdentityClient
	teslaOracle        TeslaClient
}

type IdentityClient interface {
	DescribeVehicle(vehicleID uint64) (*identity.VehicleDescription, error)
}

type DeviceActivityClient interface {
	DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*ch.Vehicle, error)
}

type DevicesClient interface {
	GetVehicleByTokenIdFast(ctx context.Context, in *pb_devices.GetVehicleByTokenIdFastRequest, opts ...grpc.CallOption) (*pb_devices.GetVehicleByTokenIdFastResponse, error)
}

type TeslaClient interface {
	GetVinByTokenId(ctx context.Context, in *pb_tesla.GetVinByTokenIdRequest, opts ...grpc.CallOption) (*pb_tesla.GetVinByTokenIdResponse, error)
}

func NewBaselineRewardService(
	settings *config.Settings,
	transferService *TransferService,
	dataService DeviceActivityClient,
	devicesClient DevicesClient,
	stakeChecker IdentityClient,
	week int,
	logger *zerolog.Logger,
	teslaOracle TeslaClient,
) *BaselineClient {
	return &BaselineClient{
		TransferService:    transferService,
		DataService:        dataService,
		DevicesClient:      devicesClient,
		ContractAddress:    common.HexToAddress(settings.IssuanceContractAddress),
		Week:               week,
		Logger:             logger,
		FirstAutomatedWeek: settings.FirstAutomatedWeek,
		IdentityClient:     stakeChecker,
		teslaOracle:        teslaOracle,
	}
}

func (t *BaselineClient) assignPoints() error {
	issuanceWeek := t.Week
	ctx := context.Background()

	weekStart := date.NumToWeekStart(issuanceWeek)
	weekEnd := date.NumToWeekEnd(issuanceWeek)

	var vinsUsed set.Set[string]

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

	activityQueryStart := time.Now()

	// These describe the active integrations for each device active this week.
	activeDevices, err := t.DataService.DescribeActiveDevices(ctx, weekStart, weekEnd)
	if err != nil {
		return err
	}

	t.Logger.Info().Msgf("Activity query took %s.", time.Since(activityQueryStart))

	lastWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek-1)).All(ctx, t.TransferService.db.DBS().Reader)
	if err != nil {
		return err
	}

	lastWeekByVehicleTokenID := make(map[int]*models.Reward)
	for _, reward := range lastWeekRewards {
		lastWeekByVehicleTokenID[reward.UserDeviceTokenID] = reward
	}

	for _, device := range activeDevices {
		logger := t.Logger.With().Int64("vehicleId", device.TokenID).Logger()

		vd, err := t.IdentityClient.DescribeVehicle(uint64(device.TokenID))
		if err != nil {
			if errors.Is(err, identity.ErrNotFound) {
				logger.Info().Msg("Vehicle was active during the week but was later deleted.")
				continue
			}
			return fmt.Errorf("failed to describe vehicle %d: %w", device.TokenID, err)
		}

		vOwner := vd.Owner

		thisWeek := &models.Reward{
			IssuanceWeekID:                 issuanceWeek,
			UserDeviceTokenID:              int(device.TokenID),
			UserEthereumAddress:            null.StringFrom(vOwner.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(vOwner.Hex()),
		}

		vin := ""

		if ad := vd.AftermarketDevice; ad != nil {
			conn, ok := constants.ConnsByMfrId[ad.Manufacturer.TokenID]
			if ok && slices.Contains(device.Sources, conn.Address.Hex()) {
				thisWeek.RewardsReceiverEthereumAddress = null.StringFrom(ad.Beneficiary.Hex())

				if vd.Owner != ad.Beneficiary {
					logger.Info().Msgf("Sending tokens to beneficiary %s.", ad.Beneficiary)
				}

				thisWeek.AftermarketTokenID = types.NewNullDecimal(decimal.New(int64(ad.TokenID), 0)) //new(decimal.Big).SetUint64(uint64(ad.TokenID)))
				thisWeek.AftermarketDevicePoints = int(conn.Points)
				thisWeek.IntegrationIds = append(thisWeek.IntegrationIds, conn.LegacyID)
			}
		}

		if sd := vd.SyntheticDevice; sd != nil {
			conn, ok := constants.ConnsByAddr[sd.Connection.Address]
			if ok && slices.Contains(device.Sources, conn.Address.Hex()) {
				thisWeek.SyntheticDeviceID = null.IntFrom(sd.TokenID)
				thisWeek.SyntheticDevicePoints = int(conn.Points)
				thisWeek.IntegrationIds = append(thisWeek.IntegrationIds, conn.LegacyID)

				// Yes, this is bad. Temporary.
				if conn.LegacyVendor == "Tesla" {
					vti, err := t.teslaOracle.GetVinByTokenId(ctx, &pb_tesla.GetVinByTokenIdRequest{TokenId: uint32(sd.TokenID)})
					if err != nil {
						if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
							// logger.Info().Msg("Device was active during the week but was later deleted.")
							continue
						}
						return err
					}
					vin = vti.Vin
				}
			}

		}

		if len(thisWeek.IntegrationIds) == 0 {
			logger.Warn().Msg("All integrations sending signals failed on-chain checks.")
			continue
		}

		if vin == "" {
			vv, err := t.DevicesClient.GetVehicleByTokenIdFast(ctx, &pb_devices.GetVehicleByTokenIdFastRequest{TokenId: uint32(device.TokenID)})
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
					// logger.Info().Msg("Device was active during the week but was later deleted.")
					continue
				}
				return err
			}
			if vv.Vin == "" {
				continue
			}
			vin = vv.Vin
		}

		if vinsUsed.Contains(vin) {
			logger.Info().Msg("VIN already used in this rewards period.")
			continue
		}

		vinsUsed.Add(vin)

		// Streak rewards.
		streakInput := StreakInput{
			ConnectedThisWeek:           true,
			ExistingConnectionStreak:    0,
			ExistingDisconnectionStreak: 0,
		}
		if lastWeek, ok := lastWeekByVehicleTokenID[int(device.TokenID)]; ok {
			streakInput.ExistingConnectionStreak = lastWeek.ConnectionStreak
			streakInput.ExistingDisconnectionStreak = lastWeek.DisconnectionStreak
		}

		streak := ComputeStreak(streakInput)

		stakePoints := 0
		if vd.Stake != nil && weekEnd.Before(vd.Stake.EndsAt) {
			stakePoints = vd.Stake.Points
			logger.Info().Msgf("Adding %d points from staking.", stakePoints)
		}

		setStreakFields(thisWeek, streak, stakePoints)

		// Anything left in this map is considered disconnected.
		// This is a no-op if the device doesn't have a record from last week.
		delete(lastWeekByVehicleTokenID, int(device.TokenID))

		// If this VIN has never earned before, make note of that.
		// Used by referrals, not this job. Have to be careful about VINs because
		// people put garbage in there.
		if len(vin) == 17 {
			vinRec := models.Vin{
				Vin:                 vin,
				FirstEarningWeek:    issuanceWeek,
				FirstEarningTokenID: types.NewDecimal(decimal.New(device.TokenID, 0)),
			}
			if err := vinRec.Upsert(ctx, t.TransferService.db.DBS().Writer, false, []string{models.VinColumns.Vin}, boil.Infer(), boil.Infer()); err != nil {
				return err
			}
		}

		if err := thisWeek.Insert(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	// We didn't see any data for these remaining devices this week.
	for _, lastWeek := range lastWeekByVehicleTokenID {
		thisWeek := &models.Reward{
			IssuanceWeekID:    issuanceWeek,
			UserDeviceTokenID: lastWeek.UserDeviceTokenID,
		}
		streakInput := StreakInput{
			ConnectedThisWeek:           false,
			ExistingConnectionStreak:    lastWeek.ConnectionStreak,
			ExistingDisconnectionStreak: lastWeek.DisconnectionStreak,
		}
		streak := ComputeStreak(streakInput)
		if streak.ConnectionStreak == 0 {
			// Don't keep these dead rows around.
			continue
		}
		setStreakFields(thisWeek, streak, 0)
		if err := thisWeek.Insert(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
			return err
		}
	}

	return nil
}

func (t *BaselineClient) calculateTokens() error {
	st := storage.DBStorage{DBS: t.TransferService.db, Logger: t.Logger}
	return st.AssignTokens(context.TODO(), t.Week, t.FirstAutomatedWeek)
}

func (t *BaselineClient) BaselineIssuance() error {
	ctx := context.Background()

	err := t.assignPoints()
	if err != nil {
		return fmt.Errorf("failed to assign points to vehicles: %w", err)
	}

	// TODO(elffjs): This blows up with a division by zero if there are no points at all.
	err = t.calculateTokens()
	if err != nil {
		return fmt.Errorf("failed to convert points into tokens: %w", err)
	}

	err = t.transferTokens(ctx)
	if err != nil {
		return fmt.Errorf("failed to submit baseline token transfers: %w", err)
	}

	week, err := models.IssuanceWeeks(models.IssuanceWeekWhere.ID.EQ(t.Week)).One(ctx, t.TransferService.db.DBS().Writer)
	if err != nil {
		return err
	}

	week.JobStatus = models.IssuanceWeeksJobStatusFinished
	if _, err := week.Update(ctx, t.TransferService.db.DBS().Writer, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func setStreakFields(reward *models.Reward, streakOutput StreakOutput, stakePoints int) {
	reward.ConnectionStreak = streakOutput.ConnectionStreak
	reward.DisconnectionStreak = streakOutput.DisconnectionStreak
	reward.StreakPoints = streakOutput.Points + stakePoints
}
