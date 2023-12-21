package services

import (
	"context"
	"github.com/DIMO-Network/rewards-api/internal/utils"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
	"testing"

	pb_defs "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/rs/zerolog"
)

type RewardsMigrationTestSuite struct {
	suite.Suite
	ctx              context.Context
	logger           zerolog.Logger
	integrationsByID map[string]*pb_defs.Integration
	allIntegrations  *pb_defs.GetIntegrationResponse
}

func (o *RewardsMigrationTestSuite) SetupSuite() {
	o.ctx = context.Background()
	o.logger = zerolog.Nop()

	fkc := &FakeDefClient{}
	o.allIntegrations, _ = fkc.GetIntegrations(o.ctx, nil, nil)

	integrsByID := make(map[string]*pb_defs.Integration)
	for _, integr := range o.allIntegrations.Integrations {
		integrsByID[integr.Id] = integr
	}
	o.integrationsByID = integrsByID
}

// TearDownTest after each test truncate tables
func (o *RewardsMigrationTestSuite) TearDownTest() {}

// TearDownSuite cleanup at end by terminating container
func (o *RewardsMigrationTestSuite) TearDownSuite() {}

// Test Runner
func TestRewardsMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(RewardsMigrationTestSuite))
}

func (o *RewardsMigrationTestSuite) TestMigrateOldRewards() {
	cont, conn := utils.GetDbConnection(o.ctx, o.T(), o.logger)
	defer func() {
		err := cont.Terminate(o.ctx)
		assert.NoError(o.T(), err)
	}()

	issuanceWeeks := []models.IssuanceWeek{
		{ID: 1, JobStatus: models.IssuanceWeeksJobStatusFinished},
	}

	for _, isw := range issuanceWeeks {
		err := isw.Insert(o.ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(o.T(), err)
	}

	existingRwrds := []models.Reward{
		{
			IssuanceWeekID:      1,
			UserDeviceID:        "MockDeviceID1",
			UserID:              "User1",
			ConnectionStreak:    13,
			DisconnectionStreak: 0,
			StreakPoints:        2000,
			IntegrationPoints:   int(o.integrationsByID[autoPiIntegration].Points) + int(o.integrationsByID[teslaIntegration].Points),
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
		err := rwrd.Insert(o.ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(o.T(), err)
	}

	err := MigrateRewardsService(o.ctx, &o.logger, conn, o.allIntegrations, 2)
	assert.NoError(o.T(), err)

	rewards, err := models.Rewards().All(o.ctx, conn.DBS().Reader)
	assert.NoError(o.T(), err)

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

	assert.Equal(o.T(), expected, actual)
}

func (o *RewardsMigrationTestSuite) Test_MultipleWeeks_MigrateOldRewards() {
	cont, conn := utils.GetDbConnection(o.ctx, o.T(), o.logger)
	defer func() {
		err := cont.Terminate(o.ctx)
		assert.NoError(o.T(), err)
	}()

	issuanceWeeks := []models.IssuanceWeek{
		{ID: 1, JobStatus: models.IssuanceWeeksJobStatusFinished},
		{ID: 2, JobStatus: models.IssuanceWeeksJobStatusFinished},
	}

	for _, isw := range issuanceWeeks {
		err := isw.Insert(o.ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(o.T(), err)
	}

	existingRwrds := []models.Reward{
		{
			IssuanceWeekID:      1,
			UserDeviceID:        "MockDeviceID1",
			UserID:              "User1",
			ConnectionStreak:    13,
			DisconnectionStreak: 0,
			StreakPoints:        2000,
			IntegrationPoints:   int(o.integrationsByID[autoPiIntegration].Points) + int(o.integrationsByID[teslaIntegration].Points),
			Tokens:              types.NewNullDecimal(decimal.New(int64(7000), 0)),
			UserEthereumAddress: null.StringFrom(mkAddr(1).Hex()),
			TransferSuccessful:  null.BoolFrom(true),
			UserDeviceTokenID:   types.NewNullDecimal(decimal.New(int64(1), 0)),
			IntegrationIds: types.StringArray{
				autoPiIntegration,
				teslaIntegration,
			},
		},
		{
			IssuanceWeekID:      2,
			UserDeviceID:        "MockDeviceID1",
			UserID:              "User1",
			ConnectionStreak:    14,
			DisconnectionStreak: 0,
			StreakPoints:        5000,
			IntegrationPoints:   int(o.integrationsByID[autoPiIntegration].Points) + int(o.integrationsByID[teslaIntegration].Points),
			Tokens:              types.NewNullDecimal(decimal.New(int64(10000), 0)),
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
		err := rwrd.Insert(o.ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(o.T(), err)
	}

	err := MigrateRewardsService(o.ctx, &o.logger, conn, o.allIntegrations, 4)
	assert.NoError(o.T(), err)

	rewards, err := models.Rewards().All(o.ctx, conn.DBS().Reader)
	assert.NoError(o.T(), err)

	expected := []models.Reward{
		{
			SyntheticDevicePoints:   4000,
			SyntheticDeviceTokens:   types.NewNullDecimal(decimal.New(int64(2333), 0)),
			AftermarketDevicePoints: 6000,
			AftermarketDeviceTokens: types.NewNullDecimal(decimal.New(int64(3500), 0)),
			StreakPoints:            2000,
			StreakTokens:            types.NewNullDecimal(decimal.New(int64(1166), 0)),
		},
		{
			SyntheticDevicePoints:   4000,
			SyntheticDeviceTokens:   types.NewNullDecimal(decimal.New(int64(2666), 0)),
			AftermarketDevicePoints: 6000,
			AftermarketDeviceTokens: types.NewNullDecimal(decimal.New(int64(4000), 0)),
			StreakPoints:            5000,
			StreakTokens:            types.NewNullDecimal(decimal.New(int64(3333), 0)),
		},
	}

	actual := []models.Reward{}

	for _, reward := range rewards {
		actual = append(actual, models.Reward{
			SyntheticDevicePoints:   reward.SyntheticDevicePoints,
			SyntheticDeviceTokens:   reward.SyntheticDeviceTokens,
			AftermarketDevicePoints: reward.AftermarketDevicePoints,
			AftermarketDeviceTokens: reward.AftermarketDeviceTokens,
			StreakPoints:            reward.StreakPoints,
			StreakTokens:            reward.StreakTokens,
		})
	}

	assert.Equal(o.T(), expected, actual)
}

func (o *RewardsMigrationTestSuite) Test_NoStreakPoints() {
	cont, conn := utils.GetDbConnection(o.ctx, o.T(), o.logger)
	defer func() {
		err := cont.Terminate(o.ctx)
		assert.NoError(o.T(), err)
	}()

	issuanceWeeks := []models.IssuanceWeek{
		{ID: 1, JobStatus: models.IssuanceWeeksJobStatusFinished},
	}

	for _, isw := range issuanceWeeks {
		err := isw.Insert(o.ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(o.T(), err)
	}

	existingRwrds := []models.Reward{
		{
			IssuanceWeekID:      1,
			UserDeviceID:        "MockDeviceID1",
			UserID:              "User1",
			ConnectionStreak:    13,
			DisconnectionStreak: 0,
			StreakPoints:        0,
			IntegrationPoints:   int(o.integrationsByID[autoPiIntegration].Points) + int(o.integrationsByID[teslaIntegration].Points),
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
		err := rwrd.Insert(o.ctx, conn.DBS().Writer, boil.Infer())
		assert.NoError(o.T(), err)
	}

	err := MigrateRewardsService(o.ctx, &o.logger, conn, o.allIntegrations, 2)
	assert.NoError(o.T(), err)

	rewards, err := models.Rewards().All(o.ctx, conn.DBS().Reader)
	assert.NoError(o.T(), err)

	expected := models.Reward{
		SyntheticDevicePoints:   4000,
		SyntheticDeviceTokens:   types.NewNullDecimal(decimal.New(int64(2800), 0)),
		AftermarketDevicePoints: 6000,
		AftermarketDeviceTokens: types.NewNullDecimal(decimal.New(int64(4200), 0)),
		StreakPoints:            0,
		StreakTokens:            types.NewNullDecimal(decimal.New(int64(0), 0)),
	}

	actual := models.Reward{
		SyntheticDevicePoints:   rewards[0].SyntheticDevicePoints,
		SyntheticDeviceTokens:   rewards[0].SyntheticDeviceTokens,
		AftermarketDevicePoints: rewards[0].AftermarketDevicePoints,
		AftermarketDeviceTokens: rewards[0].AftermarketDeviceTokens,
		StreakPoints:            rewards[0].StreakPoints,
		StreakTokens:            rewards[0].StreakTokens,
	}

	assert.Equal(o.T(), expected, actual)
}
