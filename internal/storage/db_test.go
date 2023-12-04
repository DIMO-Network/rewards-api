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
		_ = cont.Terminate(ctx)
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
		IssuanceWeekID:          80,
		UserDeviceID:            userDeviceID1,
		ConnectionStreak:        6,
		StreakPoints:            1000,
		AftermarketDevicePoints: 0,
		SyntheticDevicePoints:   1000,
	}

	reward2 := models.Reward{
		IssuanceWeekID:          80,
		UserDeviceID:            userDeviceID2,
		ConnectionStreak:        2,
		StreakPoints:            0,
		AftermarketDevicePoints: 0,
		SyntheticDevicePoints:   4000,
	}

	require.NoError(t, reward1.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer()))
	require.NoError(t, reward2.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer()))

	db := DBStorage{DBS: conn, Logger: &logger}
	err = db.AssignTokens(context.TODO(), 80, 40)
	require.NoError(t, err)

	r, _ := models.Rewards().All(context.TODO(), conn.DBS().Reader)

	fmt.Println(r)

	require.NoError(t, reward1.Reload(context.TODO(), conn.DBS().Reader))
	require.NoError(t, reward2.Reload(context.TODO(), conn.DBS().Reader))

	expect1, _ := new(big.Int).SetString("184166666666666666666666", 10)
	expect2, _ := new(big.Int).SetString("736666666666666666666666", 10)

	assert.Equal(t, expect1, reward1.SyntheticDeviceTokens.Int(nil))
	assert.Equal(t, expect2, reward2.SyntheticDeviceTokens.Int(nil))
}

func TestTokenAssignmentOneDecrease(t *testing.T) {
	ctx := context.Background()

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer func() {
		_ = cont.Terminate(ctx)
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
		IssuanceWeekID:          92,
		UserDeviceID:            userDeviceID1,
		ConnectionStreak:        6,
		StreakPoints:            1000,
		SyntheticDevicePoints:   0,
		AftermarketDevicePoints: 1000,
	}

	reward2 := models.Reward{
		IssuanceWeekID:          92,
		UserDeviceID:            userDeviceID2,
		ConnectionStreak:        2,
		StreakPoints:            0,
		SyntheticDevicePoints:   0,
		AftermarketDevicePoints: 4000,
	}

	require.NoError(t, reward1.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer()))
	require.NoError(t, reward2.Insert(context.TODO(), conn.DBS().Writer.DB, boil.Infer()))

	db := DBStorage{DBS: conn, Logger: &logger}
	err = db.AssignTokens(context.TODO(), 92, 40)
	require.NoError(t, err)

	require.NoError(t, reward1.Reload(context.TODO(), conn.DBS().Reader))
	require.NoError(t, reward2.Reload(context.TODO(), conn.DBS().Reader))

	expect1, _ := new(big.Int).SetString("156541666666666666666666", 10)
	expect2, _ := new(big.Int).SetString("626166666666666666666666", 10)

	assert.Equal(t, expect1, reward1.AftermarketDeviceTokens.Int(nil))
	assert.Equal(t, expect2, reward2.AftermarketDeviceTokens.Int(nil))
}
