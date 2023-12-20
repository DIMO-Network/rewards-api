package services

import (
	"context"
	"testing"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func TestMigrateOldRewards(t *testing.T) {
	ctx := context.Background()

	logger := zerolog.Nop()

	cont, conn := utils.GetDbConnection(ctx, t, logger)
	defer func() {
		err := cont.Terminate(ctx)
		assert.NoError(t, err)
	}()

	fkc := &FakeDefClient{}
	intgs, err := fkc.GetIntegrations(ctx, nil, nil)
	assert.NoError(t, err)

	integrsByID := make(map[string]*pb_defs.Integration)
	for _, integr := range intgs.Integrations {
		integrsByID[integr.Id] = integr
	}

	issuanceWeeks := []models.IssuanceWeek{
		{ID: 1, JobStatus: models.IssuanceWeeksJobStatusFinished},
	}

	for _, isw := range issuanceWeeks {
		err := isw.Insert(ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(t, err)
	}

	existingRwrds := []models.Reward{
		{
			IssuanceWeekID:      1,
			UserDeviceID:        "MockDeviceID1",
			UserID:              "User1",
			ConnectionStreak:    13,
			DisconnectionStreak: 0,
			StreakPoints:        2000,
			IntegrationPoints:   int(integrsByID[autoPiIntegration].Points) + int(integrsByID[teslaIntegration].Points),
			Tokens:              types.NewNullDecimal(decimal.New(int64(7000), 0)),
			UserEthereumAddress: null.StringFrom(mkAddr(1).Hex()),
			TransferSuccessful:  null.BoolFrom(true),
			UserDeviceTokenID:   types.NewNullDecimal(decimal.New(int64(1), 0)),
			IntegrationIds: types.StringArray{
				autoPiIntegration,
				teslaIntegration,
			},
		},
	}

	for _, rwrd := range existingRwrds {
		err := rwrd.Insert(ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(t, err)
	}

	err = MigrateRewardsService(ctx, &logger, conn, intgs, 2)
	assert.NoError(t, err)

	rewards, err := models.Rewards().All(ctx, conn.DBS().Reader)
	assert.NoError(t, err)

	expected := models.Reward{
		SyntheticDevicePoints:   4000,
		SyntheticDeviceTokens:   types.NewNullDecimal(decimal.New(int64(2333), 0)),
		AftermarketDevicePoints: 6000,
		AftermarketDeviceTokens: types.NewNullDecimal(decimal.New(int64(3500), 0)),
		StreakPoints:            2000,
		StreakTokens:            types.NewNullDecimal(decimal.New(int64(1166), 0)),
	}

	actual := models.Reward{
		SyntheticDevicePoints:   rewards[0].SyntheticDevicePoints,
		SyntheticDeviceTokens:   rewards[0].SyntheticDeviceTokens,
		AftermarketDevicePoints: rewards[0].AftermarketDevicePoints,
		AftermarketDeviceTokens: rewards[0].AftermarketDeviceTokens,
		StreakPoints:            rewards[0].StreakPoints,
		StreakTokens:            rewards[0].StreakTokens,
	}

	assert.Equal(t, expected, actual)
}
