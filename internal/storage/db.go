package storage

import (
	"context"
	"math/big"

	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var initialWeeklyTokens = new(big.Int).Mul(big.NewInt(1_105_000), ether)

type DBStorage struct {
	DBS    db.Store
	Logger *zerolog.Logger
}

// $1 is the week, $2 is the value to be disbursed.
const assignTokensQuery = `
UPDATE
	rewards_api.rewards
SET
  streak_tokens =
    div(
      streak_points * $2::numeric,
      (SELECT sum(streak_points + synthetic_device_points + aftermarket_device_points) FROM rewards_api.rewards WHERE issuance_week_id = $1)
    ),
  synthetic_device_tokens = 
    div(
      synthetic_device_points * $2::numeric,
      (SELECT sum(streak_points + synthetic_device_points + aftermarket_device_points) FROM rewards_api.rewards WHERE issuance_week_id = $1)
    ),
  aftermarket_device_tokens = 
    div(
      aftermarket_device_points * $2::numeric,
      (SELECT sum(streak_points + synthetic_device_points + aftermarket_device_points) FROM rewards_api.rewards WHERE issuance_week_id = $1)
    )
WHERE
	issuance_week_id = $1;`

func (s *DBStorage) AssignTokens(ctx context.Context, issuanceWeek, firstAutomatedWeek int) error {
	contractWeek := issuanceWeek - firstAutomatedWeek
	contractYear := contractWeek / 52 // This is how many years the contract thinks have passed.

	weekLimit := initialWeeklyTokens
	for i := 0; i < contractYear; i++ {
		weekLimit = new(big.Int).Mul(weekLimit, big.NewInt(85))
		weekLimit = new(big.Int).Quo(weekLimit, big.NewInt(100))
	}

	s.Logger.Info().Msgf("Database week %d, contract week %d, so contract year %d, so distributing %d tokens.", issuanceWeek, contractWeek, contractYear, weekLimit)

	_, err := s.DBS.DBS().Writer.ExecContext(ctx, assignTokensQuery, issuanceWeek, weekLimit.String())
	return err
}
