package services

import (
	"context"
	"errors"
	"fmt"

	pbdefs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
)

func MigrateRewardsService(ctx context.Context, logger *zerolog.Logger, pdb db.Store, allIntegrations *pbdefs.GetIntegrationResponse, week int) error {
	integrsByID := make(map[string]*pbdefs.Integration)
	for _, integr := range allIntegrations.Integrations {
		integrsByID[integr.Id] = integr
	}

	rewards, err := models.Rewards(
		models.RewardWhere.IssuanceWeekID.LTE(week),
		models.RewardWhere.Tokens.GT(types.NewNullDecimal(decimal.New(0, 0))),
	).All(ctx, pdb.DBS().Reader)
	if err != nil {
		return errors.New("failed to fetch rewards")
	}

	for _, reward := range rewards {
		migLogger := logger.With().Int("issuanceWeek", reward.IssuanceWeekID).Str("userDeviceID", reward.UserDeviceID).Logger()

		migLogger.Info().Msg("Starting reward token migration")

		if len(reward.IntegrationIds) == 0 {
			migLogger.Warn().Msg("skipping, could not find any integrations for device")
			continue
		}

		adPoints, sdPoints := int64(0), int64(0)

		for _, id := range reward.IntegrationIds {
			integration := integrsByID[id]
			if integration.ManufacturerTokenId == 0 { // Synthetic
				sdPoints = integration.Points
			} else {
				adPoints = integration.Points
			}
		}
		qry := `UPDATE
		rewards_api.rewards
		SET
			synthetic_device_points = $1::integer,
			aftermarket_device_points = $2::integer,
			synthetic_device_tokens = div($1 * tokens, (streak_points + integration_points)),
			aftermarket_device_tokens = div($2 * tokens, (streak_points + integration_points)),
			streak_tokens = div(streak_points * tokens, (streak_points + integration_points))
		WHERE
			issuance_week_id = $3
			AND user_device_id = $4;`

		_, err = pdb.DBS().Writer.ExecContext(ctx, qry, sdPoints, adPoints, reward.IssuanceWeekID, reward.UserDeviceID)
		if err != nil {
			migLogger.Info().Msg("Error occurred migrating rewards")
			return fmt.Errorf("error occurred splitting up rewards with issuance week %d", reward.IssuanceWeekID)
		}
	}

	return nil
}
