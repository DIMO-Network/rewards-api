package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
)

var (
	ether               = big.NewInt(params.Ether)
	etherDecimal        = new(decimal.Big).SetBigMantScale(ether, 0)
	initialWeeklyTokens = new(big.Int).Mul(big.NewInt(1_105_000), ether)
)

type DBStorage struct {
	DBS    db.Store
	Logger *zerolog.Logger
}

// $1 is the week, $2 is the value to be disbursed.
const assignTokensQuery = `UPDATE
    rewards_api.rewards
SET
    streak_tokens = div(streak_points * $2::numeric, (
            SELECT
                sum(streak_points + synthetic_device_points + aftermarket_device_points)
            FROM rewards_api.rewards
            WHERE
                issuance_week_id = $1)),
    synthetic_device_tokens = div(synthetic_device_points * $2::numeric, (
            SELECT
                sum(streak_points + synthetic_device_points + aftermarket_device_points)
            FROM rewards_api.rewards
            WHERE
                issuance_week_id = $1)),
    aftermarket_device_tokens = div(aftermarket_device_points * $2::numeric, (
            SELECT
                sum(streak_points + synthetic_device_points + aftermarket_device_points)
            FROM rewards_api.rewards
            WHERE
                issuance_week_id = $1))
WHERE
    issuance_week_id = $1;`

func (s *DBStorage) AssignTokens(ctx context.Context, issuanceWeek, firstAutomatedWeek int) error {
	// check if issuance week has any points to avoid divide by zero
	type RwrdInfo struct {
		PointsSum int `boil:"points_sum"`
	}
	var rwrdInfo RwrdInfo
	err := models.NewQuery(
		qm.Select("sum(streak_points + synthetic_device_points + aftermarket_device_points) as points_sum"),
		qm.From(models.TableNames.Rewards),
		qm.Where(models.RewardColumns.IssuanceWeekID+"=$1", issuanceWeek),
	).Bind(ctx, s.DBS.DBS().Reader, &rwrdInfo)
	if err != nil {
		return err
	}

	if rwrdInfo.PointsSum == 0 {
		return nil
	}

	contractWeek := issuanceWeek - firstAutomatedWeek
	contractYear := contractWeek / 52 // This is how many years the contract thinks have passed.

	weekLimit := initialWeeklyTokens
	for i := 0; i < contractYear; i++ {
		weekLimit = new(big.Int).Mul(weekLimit, big.NewInt(85))
		weekLimit = new(big.Int).Quo(weekLimit, big.NewInt(100))
	}

	s.Logger.Info().Msgf("Database week %d, contract week %d, so contract year %d, so distributing %d tokens.", issuanceWeek, contractWeek, contractYear, weekLimit)

	_, err = s.DBS.DBS().Writer.ExecContext(ctx, assignTokensQuery, issuanceWeek, weekLimit.String())
	return err
}

var tokensPerWeekQuery = `
SELECT 
	SUM(
		r.streak_tokens +
		r.synthetic_device_tokens +
		r.aftermarket_device_tokens
	)::numeric / NULLIF(SUM(
        r.streak_points + 
        r.aftermarket_device_points + 
        r.synthetic_device_points
    ), 0) * $2 as tokens_for_points
FROM rewards_api.rewards r 
WHERE r.issuance_week_id = $1
	AND (r.streak_tokens > 0
	OR r.synthetic_device_tokens > 0
	OR r.aftermarket_device_tokens > 0
)
LIMIT 10;
`

// CalculateTokensForPoints calculates how many tokens a given number of points is worth.
func CalculateTokensForPoints(ctx context.Context, dbStore db.Store, points int, weekID int) (*decimal.Big, error) {
	var tokens types.NullDecimal
	err := dbStore.DBS().Reader.QueryRowContext(ctx, tokensPerWeekQuery, weekID, points).Scan(&tokens)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no conversion rate found for weekId %v", weekID)
		}
		return nil, fmt.Errorf("error calculating tokens: %w", err)
	}
	if tokens.Big == nil {
		return nil, fmt.Errorf("null conversion rate found for weekId %v", weekID)
	}

	// Divide result by ether
	result := new(decimal.Big).Quo(tokens.Big, etherDecimal)
	return result, nil
}
