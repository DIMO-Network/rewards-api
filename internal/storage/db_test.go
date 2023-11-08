package storage

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestTokenAssignmentNoDecrease(t *testing.T) {
	ctx := context.Background()

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer func() {
		cont.Terminate(ctx)
	}()

	userDeviceID1 := ksuid.New().String()
	userDeviceID2 := ksuid.New().String()

	wk := models.IssuanceWeek{
		ID:        80,
		JobStatus: models.IssuanceWeeksJobStatusFinished,
	}

	err := wk.Insert(context.TODO(), conn.DBS().Writer, boil.Infer())
	require.NoError(t, err)

	reward1 := models.Reward{
		IssuanceWeekID:    80,
		UserDeviceID:      userDeviceID1,
		ConnectionStreak:  6,
		StreakPoints:      1000,
		IntegrationPoints: 1000,
	}

	reward2 := models.Reward{
		IssuanceWeekID:    80,
		UserDeviceID:      userDeviceID2,
		ConnectionStreak:  2,
		StreakPoints:      0,
		IntegrationPoints: 4000,
	}

	err = reward1.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer())
	require.NoError(t, err)
	err = reward2.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer())
	require.NoError(t, err)

	db := DBStorage{DBS: conn}
	err = db.AssignTokens(context.TODO(), 80, 40)
	require.NoError(t, err)

	r, _ := models.Rewards().All(context.TODO(), conn.DBS().Reader)

	fmt.Println(r)

	reward1.Reload(context.TODO(), conn.DBS().Reader)
	reward2.Reload(context.TODO(), conn.DBS().Reader)

	expect1, _ := new(big.Int).SetString("368333333333333333333333", 10)
	expect2, _ := new(big.Int).SetString("736666666666666666666666", 10)

	assert.Equal(t, expect1, reward1.Tokens.Int(nil))
	assert.Equal(t, reward2.Tokens.Int(nil), expect2)
}

func TestTokenAssignmentOneDecrease(t *testing.T) {
	ctx := context.Background()

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer func() {
		cont.Terminate(ctx)
	}()

	userDeviceID1 := ksuid.New().String()
	userDeviceID2 := ksuid.New().String()

	wk := models.IssuanceWeek{
		ID:        92,
		JobStatus: models.IssuanceWeeksJobStatusFinished,
	}

	err := wk.Insert(context.TODO(), conn.DBS().Writer, boil.Infer())
	require.NoError(t, err)

	reward1 := models.Reward{
		IssuanceWeekID:    92,
		UserDeviceID:      userDeviceID1,
		ConnectionStreak:  6,
		StreakPoints:      1000,
		IntegrationPoints: 1000,
	}

	reward2 := models.Reward{
		IssuanceWeekID:    92,
		UserDeviceID:      userDeviceID2,
		ConnectionStreak:  2,
		StreakPoints:      0,
		IntegrationPoints: 4000,
	}

	err = reward1.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer())
	require.NoError(t, err)
	err = reward2.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer())
	require.NoError(t, err)

	db := DBStorage{DBS: conn}
	err = db.AssignTokens(context.TODO(), 92, 40)
	require.NoError(t, err)

	r, _ := models.Rewards().All(context.TODO(), conn.DBS().Reader)

	fmt.Println(r)

	reward1.Reload(context.TODO(), conn.DBS().Reader)
	reward2.Reload(context.TODO(), conn.DBS().Reader)

	expect1, _ := new(big.Int).SetString("313083333333333333333333", 10)
	expect2, _ := new(big.Int).SetString("626166666666666666666666", 10)

	assert.Equal(t, reward1.Tokens.Int(nil), expect1)
	assert.Equal(t, reward2.Tokens.Int(nil), expect2)
}
