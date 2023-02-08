package storage

import (
	"context"
	"math/big"

	"github.com/DIMO-Network/shared/db"
)

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

func (s *DBStorage) AssignTokens(ctx context.Context, issuanceWeek int, totalTokens *big.Int) error {
	_, err := s.DBS.DBS().Writer.ExecContext(ctx, assignTokensQuery, issuanceWeek, totalTokens.String())
	return err
}
