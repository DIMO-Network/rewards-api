package storage

import (
	"context"
	"math/big"

	"github.com/DIMO-Network/shared/db"
)

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var initialWeeklyTokens = new(big.Int).Mul(big.NewInt(1_105_000), ether)

type DBStorage struct {
	DBS db.Store
}

// $1 is the week, $2 is the value to be disbursed.
const assignTokensQuery = `
UPDATE
	rewards_api.rewards
SET
	tokens =
		div(
			(streak_points + integration_points) * $2::numeric,
			(SELECT sum(streak_points + integration_points) FROM rewards_api.rewards WHERE issuance_week_id = $1)
		)
WHERE
	issuance_week_id = $1;`

func (s *DBStorage) AssignTokens(ctx context.Context, issuanceWeek, firstAutomatedWeek int) error {
	// This is how many years the contract thinks have passed.
	contractYear := (issuanceWeek - firstAutomatedWeek) / 52

	weekLimit := initialWeeklyTokens
	for i := 0; i < contractYear; i++ {
		weekLimit = new(big.Int).Mul(weekLimit, big.NewInt(85))
		weekLimit = new(big.Int).Quo(weekLimit, big.NewInt(100))
	}

	_, err := s.DBS.DBS().Writer.ExecContext(ctx, assignTokensQuery, issuanceWeek, weekLimit.String())
	return err
}
