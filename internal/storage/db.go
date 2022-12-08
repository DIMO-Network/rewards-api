package storage

import (
	"context"
	"math/big"

	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/DIMO-Network/rewards-api/models"
)

type Storage interface {
	AssignTokens(ctx context.Context, issuanceWeek int, totalTokens *big.Int) error
	Rewards(ctx context.Context, issuanceWeek int) ([]*Reward, error)
}

func NewDB(db func() *database.DBReaderWriter) Storage {
	return &dbStorage{db: db}
}

type dbStorage struct {
	db func() *database.DBReaderWriter
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

	_, err := s.db().Writer.ExecContext(ctx, q, issuanceWeek, totalTokens.String())
	return err
}

type Reward struct {
	UserDeviceID string
	Tokens       *big.Int
}

func (s *dbStorage) Rewards(ctx context.Context, issuanceWeek int) ([]*Reward, error) {
	rows, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).All(ctx, s.db().Reader)
	if err != nil {
		return nil, err
	}
	out := make([]*Reward, len(rows))

	for i, row := range rows {
		out[i] = &Reward{UserDeviceID: row.UserDeviceID, Tokens: row.Tokens.Int(nil)}
	}

	return out, nil
}
