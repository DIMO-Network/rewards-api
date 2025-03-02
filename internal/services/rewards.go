//go:generate mockgen -source=./rewards.go -destination=rewards_mock_test.go -package=services
package services

import (
	"context"
	"fmt"
	"time"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	pb_devices "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/date"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/storage"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/set"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/volatiletech/null/v8"

	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BaselineClient struct {
	TransferService    *TransferService
	DataService        DeviceActivityClient
	DevicesClient      DevicesClient
	DefsClient         IntegrationsGetter
	vinVCSrv           VINVCService
	ContractAddress    common.Address
	Week               int
	Logger             *zerolog.Logger
	FirstAutomatedWeek int
	StakeChecker       StakeChecker
	StakingEnabled     bool
}

type StakeChecker interface {
	GetVehicleStakePoints(vehicleID uint64) (int, error)
}

type VINVCService interface {
	GetConfirmedVINVCs(ctx context.Context, activeTokenIds []*ch.Vehicle) (map[int64]struct{}, error)
}
type DeviceActivityClient interface {
	DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*ch.Vehicle, error)
}

type IntegrationsGetter interface {
	GetIntegrations(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb_defs.GetIntegrationResponse, error)
}

type DevicesClient interface {
	GetUserDeviceByTokenId(ctx context.Context, in *pb_devices.GetUserDeviceByTokenIdRequest, opts ...grpc.CallOption) (*pb_devices.UserDevice, error)
}

func NewBaselineRewardService(
	settings *config.Settings,
	transferService *TransferService,
	dataService DeviceActivityClient,
	devicesClient DevicesClient,
	defsClient IntegrationsGetter,
	stakeChecker StakeChecker,
	vinVCSrv VINVCService,
	week int,
	logger *zerolog.Logger,
) *BaselineClient {
	return &BaselineClient{
		TransferService:    transferService,
		DataService:        dataService,
		DevicesClient:      devicesClient,
		DefsClient:         defsClient,
		ContractAddress:    common.HexToAddress(settings.IssuanceContractAddress),
		Week:               week,
		Logger:             logger,
		FirstAutomatedWeek: settings.FirstAutomatedWeek,
		vinVCSrv:           vinVCSrv,
		StakeChecker:       stakeChecker,
		StakingEnabled:     settings.EnableStaking,
	}
}

func (t *BaselineClient) assignPoints() error {
	issuanceWeek := t.Week
	ctx := context.Background()

	weekStart := date.NumToWeekStart(issuanceWeek)
	weekEnd := date.NumToWeekEnd(issuanceWeek)

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

	// These describe the active integrations for each device active this week.
	activeDevices, err := t.DataService.DescribeActiveDevices(ctx, weekStart, weekEnd)
	if err != nil {
		return err
	}

	vinVCConfirmed, err := t.vinVCSrv.GetConfirmedVINVCs(ctx, activeDevices)
	if err != nil {
		// this is a non-fatal error, we can continue without this data
		t.Logger.Warn().Err(err).Msg("Failed to get confirmed VIN VC VINs. continuing execution.")
		vinVCConfirmed = map[int64]struct{}{}
	}

	allIntegrations, err := t.DefsClient.GetIntegrations(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	amMfrTokenToIntegration := make(map[uint64]*pb_defs.Integration)
	swIntegrsByTokenID := make(map[uint64]*pb_defs.Integration)

	for _, integr := range allIntegrations.Integrations {
		if integr.ManufacturerTokenId == 0 {
			// Must be a software integration.
			swIntegrsByTokenID[integr.TokenId] = integr
		} else {
			// Must be the integration associated with a manufacturer.
			amMfrTokenToIntegration[integr.ManufacturerTokenId] = integr
		}
	}

	lastWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek-1)).All(ctx, t.TransferService.db.DBS().Reader)
	if err != nil {
		return err
	}

	lastWeekByDevice := make(map[string]*models.Reward)
	for _, reward := range lastWeekRewards {
		lastWeekByDevice[reward.UserDeviceID] = reward
	}

	for _, device := range activeDevices {
		logger := t.Logger.With().Int64("vehicleTokenId", device.TokenID).Logger()

		ud, err := t.DevicesClient.GetUserDeviceByTokenId(ctx, &pb_devices.GetUserDeviceByTokenIdRequest{TokenId: device.TokenID})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				logger.Info().Msg("Device was active during the week but was later deleted.")
				continue
			}
			return err
		}

		integsSignalsThisWeek := set.New(device.Integrations...)

		logger = logger.With().Str("userId", ud.UserId).Logger()

		if _, ok := vinVCConfirmed[device.TokenID]; !ok && len(vinVCConfirmed) > 0 {
			// TODO: Update this to a continue after we have a better idea of how many vehicles are missing VIN VC.
			logger.Warn().Str("deviceId", ud.Id).Bool("vinConfirmed", ud.VinConfirmed).Msg("Vehicle does not have a confirmed VIN VC VIN.")
		}

		if !ud.VinConfirmed {
			logger.Info().Msg("Device does not have confirmed VIN.")
			continue
		}

		if len(ud.OwnerAddress) != 20 {
			logger.Info().Msg("User has minted a car but has no owner address?")
			continue
		}

		vOwner := common.BytesToAddress(ud.OwnerAddress)

		thisWeek := &models.Reward{
			UserDeviceID:                   ud.Id,
			IssuanceWeekID:                 issuanceWeek,
			UserID:                         ud.UserId,
			UserDeviceTokenID:              types.NewNullDecimal(new(decimal.Big).SetUint64(*ud.TokenId)),
			UserEthereumAddress:            null.StringFrom(vOwner.Hex()),
			RewardsReceiverEthereumAddress: null.StringFrom(vOwner.Hex()),
		}

		if ad := ud.AftermarketDevice; ad != nil {
			// Want to see if this kind (right manufacturer) of device transmitted for this vehicle
			// this week.
			if ad.ManufacturerTokenId == 0 {
				logger.Warn().Msgf("Aftermarket device %d does not have a manufacturer.", ad.TokenId)
				continue
			}

			integr, ok := amMfrTokenToIntegration[ad.ManufacturerTokenId]
			if !ok {
				logger.Warn().Msgf("Aftermarket device manufacturer %d does not have an associated integration.", ad.ManufacturerTokenId)
				continue
			}

			if integsSignalsThisWeek.Contains(integr.Id) {
				if len(ad.Beneficiary) == 20 {
					bene := common.BytesToAddress(ad.Beneficiary)
					if vOwner != bene {
						logger.Info().Msgf("Sending tokens to beneficiary %s for aftermarket device %d.", bene.Hex(), ad.TokenId)
						thisWeek.RewardsReceiverEthereumAddress = null.StringFrom(bene.Hex())
					}
				} else {
					logger.Warn().Msgf("Aftermarket device %d is not returning a beneficiary.", ad.TokenId)
				}

				thisWeek.AftermarketTokenID = types.NewNullDecimal(new(decimal.Big).SetUint64(ad.TokenId))
				thisWeek.AftermarketDevicePoints = int(integr.Points)
				thisWeek.IntegrationIds = append(thisWeek.IntegrationIds, integr.Id)
			}
		}

		if sd := ud.SyntheticDevice; sd != nil {
			if sd.IntegrationTokenId == 0 {
				logger.Warn().Msgf("Synthetic device %d does not have an integration.", sd.IntegrationTokenId)
				continue
			}

			integr, ok := swIntegrsByTokenID[sd.IntegrationTokenId]
			if !ok {
				logger.Warn().Msgf("Synthetic device %d has integration %d without metadata.", sd.TokenId, sd.IntegrationTokenId)
				continue
			}

			if integsSignalsThisWeek.Contains(integr.Id) {
				thisWeek.SyntheticDeviceID = null.IntFrom(int(sd.TokenId))
				thisWeek.SyntheticDevicePoints = int(integr.Points)
				thisWeek.IntegrationIds = append(thisWeek.IntegrationIds, integr.Id)
			}
		}

		if len(thisWeek.IntegrationIds) == 0 {
			logger.Warn().Msg("All integrations sending signals failed on-chain checks.")
			continue
		}

		// Streak rewards.
		streakInput := StreakInput{
			ConnectedThisWeek:           true,
			ExistingConnectionStreak:    0,
			ExistingDisconnectionStreak: 0,
		}
		if lastWeek, ok := lastWeekByDevice[ud.Id]; ok {
			streakInput.ExistingConnectionStreak = lastWeek.ConnectionStreak
			streakInput.ExistingDisconnectionStreak = lastWeek.DisconnectionStreak
		}

		streak := ComputeStreak(streakInput)

		stakePoints := 0

		if t.StakingEnabled {
			var err error
			stakePoints, err = t.StakeChecker.GetVehicleStakePoints(*ud.TokenId)
			if err != nil {
				return fmt.Errorf("failed to check staking for vehicle %d: %w", device.TokenID, err)
			}
			if stakePoints != 0 {
				logger.Info().Uint64("vehicleId", *ud.TokenId).Msgf("Adding %d points from staking.", stakePoints)
			}
		}

		setStreakFields(thisWeek, streak, stakePoints)

		// Anything left in this map is considered disconnected.
		// This is a no-op if the device doesn't have a record from last week.
		delete(lastWeekByDevice, ud.Id)

		// If this VIN has never earned before, make note of that.
		// Used by referrals, not this job. Have to be careful about VINs because
		// people put garbage in there.
		if ud.Vin != nil && len(*ud.Vin) == 17 {
			vinRec := models.Vin{
				Vin:                 *ud.Vin,
				FirstEarningWeek:    issuanceWeek,
				FirstEarningTokenID: types.NewDecimal(new(decimal.Big).SetUint64(*ud.TokenId)),
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
