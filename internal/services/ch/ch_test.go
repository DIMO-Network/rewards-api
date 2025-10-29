package ch

import (
	"context"
	"testing"
	"time"

	chconfig "github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
	"github.com/DIMO-Network/clickhouse-infra/pkg/container"
	"github.com/DIMO-Network/model-garage/pkg/migrations"
	"github.com/DIMO-Network/model-garage/pkg/vss"
	"github.com/stretchr/testify/suite"
)

type CHTestSuite struct {
	suite.Suite
	chClient  *Client
	container *container.Container
}

func TestCHService(t *testing.T) {
	suite.Run(t, new(CHTestSuite))
}

const day = 24 * time.Hour

const (
	macaron = "0x4c674ddE8189aEF6e3b58F5a36d7438b2b1f6Bc2"
	tesla   = "0xc4035Fecb1cc906130423EF05f9C20977F643722"
	autoPi  = "0x5e31bBc786D7bEd95216383787deA1ab0f1c1897"
	ruptela = "0xF26421509Efe92861a587482100c6d728aBf1CD0"
)

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

	mustAppend(&vss.Signal{TokenID: 3, Timestamp: weekEnd.Add(-10 * day), Name: "xdd", Source: ruptela, ValueNumber: 10.2})
	mustAppend(&vss.Signal{TokenID: 3, Timestamp: weekEnd.Add(-11 * day), Name: "powertrainTractionBatteryGrossCapacity", Source: tesla, ValueNumber: 65})
	mustAppend(&vss.Signal{TokenID: 3, Timestamp: weekEnd.Add(-day), Name: "xdd", Source: ruptela, ValueNumber: 10.2})
	mustAppend(&vss.Signal{TokenID: 3, Timestamp: weekEnd.Add(-2 * day), Name: "xdd2", Source: ruptela, ValueNumber: 10.55})
	mustAppend(&vss.Signal{TokenID: 5, Timestamp: weekEnd.Add(-3 * day), Name: "xdd3", Source: tesla, ValueNumber: 10.55})
	mustAppend(&vss.Signal{TokenID: 5, Timestamp: weekEnd.Add(-4 * day), Name: "xdd3", Source: ruptela, ValueNumber: 11.55})
	mustAppend(&vss.Signal{TokenID: 7, Timestamp: weekEnd.Add(-10 * day), Name: "xdd3", Source: tesla, ValueNumber: 10.55})
	mustAppend(&vss.Signal{TokenID: 11, Timestamp: weekEnd.Add(-4 * day), Name: "xdd", Source: macaron, ValueNumber: 7.5})
	mustAppend(&vss.Signal{TokenID: 13, Timestamp: weekEnd.Add(-3 * day), Name: "powertrainTractionBatteryGrossCapacity", Source: ruptela, ValueNumber: 75.0})

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
		vehicleIDToIntegrations[r.TokenID] = r.Sources
	}

	c.Require().ElementsMatch([]string{ruptela}, vehicleIDToIntegrations[3], "Mismatch for vehicle 3.")
	c.Require().ElementsMatch([]string{tesla, ruptela}, vehicleIDToIntegrations[5], "Mismatch for vehicle 5.")
	c.Require().ElementsMatch([]string{macaron}, vehicleIDToIntegrations[11], "Mismatch for vehicle 11.")
}

func (c *CHTestSuite) Test_GetIntegrations() {
	resp, err := c.chClient.GetSourcesForVehicles(context.TODO(), []uint64{3, 7, 11, 13}, weekEnd.Add(-7*day), weekEnd)
	c.Require().NoError(err)

	c.Len(resp, 2)

	vehicleIDToIntegrations := make(map[int64][]string)
	for _, r := range resp {
		vehicleIDToIntegrations[r.TokenID] = r.Sources
	}

	c.Require().ElementsMatch([]string{ruptela}, vehicleIDToIntegrations[3], "Mismatch for vehicle 3.")
	c.Require().ElementsMatch([]string{macaron}, vehicleIDToIntegrations[11], "Mismatch for vehicle 5.")
}
