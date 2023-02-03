package storage

import (
	"context"
	"math/big"

	"github.com/DIMO-Network/shared/db"
)

type Storage interface {
	AssignTokens(ctx context.Context, issuanceWeek int, totalTokens *big.Int) error
}

func NewDB(base db.Store) Storage {
	return &dbStorage{db: base}
}

type dbStorage struct {
	db db.Store
}

func (s *dbStorage) AssignTokens(ctx context.Context, issuanceWeek int, totalTokens *big.Int) error {
	q := `
		UPDATE rewards_api.rewards
		SET tokens =
			div(
				(streak_points + integration_points) * $2::numeric,
				(SELECT sum(streak_points + integration_points) FROM rewards_api.rewards WHERE issuance_week_id = $1)
			)
		WHERE issuance_week_id = $1;`

	_, err := s.db.DBS().Writer.ExecContext(ctx, q, issuanceWeek, totalTokens.String())
	return err
}

type Reward struct {
	UserDeviceID string
	Tokens       *big.Int
}
