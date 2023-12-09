package services

import (
	"context"
	"errors"
	"fmt"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func MigrateRewardsController(ctx context.Context, logger *zerolog.Logger, settings *config.Settings, pdb db.Store, allIntegrations *pb_defs.GetIntegrationResponse, week int) error {

	integrsByID := make(map[string]*pb_defs.Integration)
	for _, integr := range allIntegrations.Integrations {
		integrsByID[integr.Id] = integr
	}

	rewards, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.LT(week),
		models.RewardWhere.Tokens.GT(types.NewNullDecimal(decimal.New(0, 0))),
	).All(ctx, pdb.DBS().Reader)
	if err != nil {
		return errors.New("failed to fetch rewards")
	}

	for _, reward := range rewards {
		logger.Info().Int("Issuance Week", reward.IssuanceWeekID).Str("DeviceID", reward.UserDeviceID).Msg("Starting reward token migration")

		if len(reward.IntegrationIds) == 0 {
			continue
		}

		var sd *pb_defs.Integration
		var ad *pb_defs.Integration

		for _, id := range reward.IntegrationIds {
			integration := integrsByID[id]

			if integration.ManufacturerTokenId == 0 { // Synthetic
				sd = integration
			} else {
				ad = integration
			}
		}

		qry := `UPDATE
			rewards_api.rewards
			SET
				synthetic_device_points = $1,
				aftermarket_device_points = $2,
			synthetic_device_tokens = div(
				$1 * tokens::INTEGER,
				(SELECT
					sum(streak_points + integration_points)
					FROM rewards_api.rewards
					WHERE
						issuance_week_id = $3
						AND user_device_id = $4
				)
			),
			aftermarket_device_tokens = div(
				$2 * tokens::INTEGER,
				(SELECT
					sum(streak_points + integration_points)
					FROM rewards_api.rewards
					WHERE
						issuance_week_id = $3
						AND user_device_id = $4
				)
			),
			streak_tokens = div(
				streak_points * tokens::INTEGER,
				(SELECT
					sum(streak_points + integration_points)
					FROM rewards_api.rewards
					WHERE
						issuance_week_id = $3
						AND user_device_id = $4
				)
			)
			`
		_, err = pdb.DBS().Writer.ExecContext(ctx, qry, int(sd.Points), int(ad.Points), reward.IssuanceWeekID, reward.UserDeviceID)
		if err != nil {
			logger.Info().Int("Issuance Week", reward.IssuanceWeekID).Str("DeviceID", reward.UserDeviceID).Err(err).Msg("Error occurred migrating rewards")
			return fmt.Errorf("error occurred splitting up rewards with issuance week %d", reward.IssuanceWeekID)
		}
	}

	return nil
}
