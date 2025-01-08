package ch

import (
	"context"
	"testing"
	"time"

	chconfig "github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
	"github.com/DIMO-Network/clickhouse-infra/pkg/container"
	"github.com/DIMO-Network/model-garage/pkg/migrations"
	"github.com/stretchr/testify/suite"
)

type CHTestSuite struct {
	suite.Suite
	chClient  *Client
	container *container.Container
}

var Integrations = struct {
	Macaron  string
	AutoPi   string
	Smartcar string
	Tesla    string
	Ruptela  string
}{
	"2ULfuC8U9dOqRshZBAi0lMM1Rrx", // Macaron
	"27qftVRWQYpVDcO5DltO5Ojbjxk", // AutoPi
	"22N2xaPOq2WW2gAHBHd0Ikn4Zob", // Smartcar
	"26A5Dk3vvvQutjSyF0Jka2DP5lg", // Tesla
	"2lcaMFuCO0HJIUfdq8o780Kx5n3", // Ruptela
}

func TestCHService(t *testing.T) {
	suite.Run(t, new(CHTestSuite))
}

const day = 24 * time.Hour

var weekEnd = time.Date(2024, 10, 14, 5, 0, 0, 0, time.UTC)

func (c *CHTestSuite) SetupSuite() {
	ctx := context.Background()
	container, err := container.CreateClickHouseContainer(ctx, chconfig.Settings{})
	c.Require().NoError(err, "Failed to create clickhouse container")

	db, err := container.GetClickhouseAsDB()
	c.Require().NoError(err, "Failed to get clickhouse connection")

	err = migrations.RunGoose(ctx, []string{"up", "-v"}, db)
	c.Require().NoError(err, "Failed to run migrations")

	conn, err := container.GetClickHouseAsConn()
	c.Require().NoError(err, "Failed to get clickhouse connection")

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO signal")
	c.Require().NoError(err, "Failed to prepare batch")

	mustAppend := func(v any) {
		err := batch.AppendStruct(v)
		c.Require().NoError(err)
	}

	mustAppend(&signalRow{TokenID: 3, Timestamp: weekEnd.Add(-10 * day), Name: "xdd", Source: "dimo/integration/" + Integrations.Smartcar, ValueNumber: 10.2})
	mustAppend(&signalRow{TokenID: 3, Timestamp: weekEnd.Add(-day), Name: "xdd", Source: "dimo/integration/" + Integrations.Smartcar, ValueNumber: 10.2})
	mustAppend(&signalRow{TokenID: 3, Timestamp: weekEnd.Add(-2 * day), Name: "xdd2", Source: "dimo/integration/" + Integrations.Macaron, ValueNumber: 10.55})
	mustAppend(&signalRow{TokenID: 5, Timestamp: weekEnd.Add(-3 * day), Name: "xdd3", Source: "dimo/integration/" + Integrations.Tesla, ValueNumber: 10.55})
	mustAppend(&signalRow{TokenID: 7, Timestamp: weekEnd.Add(-10 * day), Name: "xdd3", Source: "dimo/integration/" + Integrations.Tesla, ValueNumber: 10.55})
	mustAppend(&signalRow{TokenID: 11, Timestamp: weekEnd.Add(-4 * day), Name: "xdd", Source: "0xF26421509Efe92861a587482100c6d728aBf1CD0", ValueNumber: 7.5})
	mustAppend(&signalRow{TokenID: 13, Timestamp: weekEnd.Add(-3 * day), Name: "powertrainTractionBatteryGrossCapacity", Source: "dimo/integration/" + Integrations.Smartcar, ValueNumber: 75.0})

	c.Require().NoError(err, "Failed to append struct")

	err = batch.Send()
	c.Require().NoError(err, "Failed to insert rows")

	err = batch.Send()
	c.Require().NoError(err, "Failed to send batch")

	c.container = container
	c.chClient = &Client{conn: conn}
}

func (c *CHTestSuite) TearDownSuite() {
	c.container.Terminate(context.Background())
}

func (c *CHTestSuite) Test_DescribeActiveDevices() {
	ctx := context.Background()

	resp, err := c.chClient.DescribeActiveDevices(ctx, weekEnd.Add(-7*day), weekEnd)
	c.Require().NoError(err)

	c.Len(resp, 3)

	vehicleIDToIntegrations := make(map[int64][]string)
	for _, r := range resp {
		vehicleIDToIntegrations[r.TokenID] = r.Integrations
	}

	c.Require().ElementsMatch(vehicleIDToIntegrations[3], []string{Integrations.Smartcar, Integrations.Macaron})
	c.Require().ElementsMatch(vehicleIDToIntegrations[5], []string{Integrations.Tesla})
	c.Require().ElementsMatch(vehicleIDToIntegrations[11], []string{Integrations.Ruptela})
}

func (c *CHTestSuite) Test_GetIntegrations() {
	resp, err := c.chClient.GetIntegrationsForVehicles(context.TODO(), []uint64{3, 7, 11, 13}, weekEnd.Add(-7*day), weekEnd)
	c.Require().NoError(err)

	c.Len(resp, 2)

	vehicleIDToIntegrations := make(map[int64][]string)
	for _, r := range resp {
		vehicleIDToIntegrations[r.TokenID] = r.Integrations
	}

	c.Require().ElementsMatch(vehicleIDToIntegrations[3], []string{Integrations.Smartcar, Integrations.Macaron})
	c.Require().ElementsMatch(vehicleIDToIntegrations[11], []string{Integrations.Ruptela})
}

type signalRow struct {
	TokenID     int64     `ch:"token_id"`
	Timestamp   time.Time `ch:"timestamp"`
	Name        string    `ch:"name"`
	Source      string    `ch:"source"`
	ValueNumber float64   `ch:"value_number"`
	ValueString string    `ch:"value_string"`
}
