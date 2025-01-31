package storage

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
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

func TestCalculateTokensForPointsPerformance(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer func() {
		_ = cont.Terminate(ctx)
	}()

	// Set a fixed conversion rate (tokens per point)
	conversionRate := decimal.New(5, 0) // 5 tokens per point

	// Create a test week
	testDate := time.Now()
	wk := models.IssuanceWeek{
		ID:        100,
		StartsAt:  testDate,
		EndsAt:    testDate.Add(7 * 24 * time.Hour),
		JobStatus: models.IssuanceWeeksJobStatusFinished,
	}

	err := wk.Insert(context.TODO(), conn.DBS().Writer, boil.Infer())
	require.NoError(t, err)

	// Insert test records
	fmt.Printf("Starting to insert test records...\n")
	insertStart := time.Now()

	for i := 0; i < 1000; i++ {
		streakPoints := rand.Int() % 1000
		aftermarketPoints := rand.Int() % 1000
		syntheticPoints := rand.Int() % 1000

		// Calculate tokens: points * conversion rate * etherx
		streakTokens := new(decimal.Big).Mul(decimal.New(int64(streakPoints), 0), conversionRate)
		streakTokens.Mul(streakTokens, etherDecimal)
		aftermarketTokens := new(decimal.Big).Mul(decimal.New(int64(aftermarketPoints), 0), conversionRate)
		aftermarketTokens.Mul(aftermarketTokens, etherDecimal)
		syntheticTokens := new(decimal.Big).Mul(decimal.New(int64(syntheticPoints), 0), conversionRate)
		syntheticTokens.Mul(syntheticTokens, etherDecimal)
		reward := models.Reward{
			IssuanceWeekID:          100,
			UserDeviceID:            ksuid.New().String(),
			UserID:                  ksuid.New().String(),
			StreakPoints:            streakPoints,
			AftermarketDevicePoints: aftermarketPoints,
			SyntheticDevicePoints:   syntheticPoints,

			CreatedAt:               testDate,
			UpdatedAt:               testDate,
			StreakTokens:            types.NewNullDecimal(streakTokens),
			AftermarketDeviceTokens: types.NewNullDecimal(aftermarketTokens),
			SyntheticDeviceTokens:   types.NewNullDecimal(syntheticTokens),
		}

		err := reward.Insert(context.TODO(), conn.DBS().Writer, boil.Infer())
		require.NoError(t, err)

	}

	insertDuration := time.Since(insertStart)
	fmt.Printf("Finished inserting records in %v\n", insertDuration)

	// Run the performance test

	startTime := time.Now()
	actualTokens, err := CalculateTokensForPoints(ctx, conn, 1000, 100)
	require.NoError(t, err)
	duration := time.Since(startTime)

	// Calculate expected result: 1000 * conversion rate
	expectedTokens := new(decimal.Big).Mul(decimal.New(1000, 0), conversionRate)

	// Convert to float64 for comparison
	expectedFloat, ok := expectedTokens.Float64()
	require.True(t, ok, "Failed to convert expected tokens to float64")
	actualFloat, ok := actualTokens.Float64()
	require.True(t, ok, "Failed to convert actual tokens to float64")

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, actualTokens)
	assert.InEpsilon(t, expectedFloat, actualFloat, 0.0001, "Token calculation outside acceptable range")

	fmt.Printf("Query execution time: %v\n", duration)

	// Performance threshold check
	assert.Less(t, duration, 1*time.Second, "Query took longer than 1 second to execute")
}
