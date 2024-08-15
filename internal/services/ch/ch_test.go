package ch

import (
	"context"
	"testing"
	"time"

	chconfig "github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
	"github.com/DIMO-Network/clickhouse-infra/pkg/container"
	dimo_ch "github.com/DIMO-Network/clickhouse-infra/pkg/container"
	"github.com/DIMO-Network/model-garage/pkg/migrations"
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

func (c *CHTestSuite) SetupSuite() {
	ctx := context.Background()
	container, err := dimo_ch.CreateClickHouseContainer(ctx, chconfig.Settings{})
	c.Require().NoError(err, "Failed to create clickhouse container")

	db, err := container.GetClickhouseAsDB()
	c.Require().NoError(err, "Failed to get clickhouse connection")

	err = migrations.RunGoose(ctx, []string{"up", "-v"}, db)
	c.Require().NoError(err, "Failed to run migrations")

	conn, err := container.GetClickHouseAsConn()
	c.Require().NoError(err, "Failed to get clickhouse connection")

	err = conn.Exec(ctx, createSignalStmt)
	c.Require().NoError(err, "Failed to insert dummy data")

	err = conn.Exec(ctx, dummyData)
	c.Require().NoError(err, "Failed to insert dummy data")

	chClient := &Client{conn: conn}

	c.chClient = chClient
	c.container = container
}

func (c *CHTestSuite) TearDownSuite() {
	c.container.Terminate(context.Background())
}

func (c *CHTestSuite) Test_DescribeActiveDevices() {
	ctx := context.Background()
	start := time.Now().AddDate(0, 0, -6)
	end := time.Now().AddDate(0, 0, 1)
	resp, err := c.chClient.DescribeActiveDevices(ctx, start, end)
	c.Require().NoError(err)
	c.Require().Equal(10, len(resp))
}

func (c *CHTestSuite) Test_GetIntegrations() {
	ctx := context.Background()
	start := time.Now().AddDate(0, 0, -6)
	end := time.Now().AddDate(0, 0, 1)
	resp, err := c.chClient.GetIntegrations(ctx, 1, start, end)
	c.Require().NoError(err)
	c.Require().Equal(3, len(resp))
}

const createSignalStmt = `
CREATE TABLE IF NOT EXISTS signal
(
	token_id UInt32 COMMENT 'token_id of this device data.',
	timestamp DateTime64(6, 'UTC') COMMENT 'timestamp of when this data was collected.',
	name LowCardinality(String) COMMENT 'name of the signal collected.',
	source String COMMENT 'source of the signal collected.',
	value_number Float64 COMMENT 'float64 value of the signal collected.',
	value_string String COMMENT 'string value of the signal collected.'
)
ENGINE = ReplacingMergeTree
ORDER BY (token_id, timestamp, name)
`

const dummyData = `
INSERT INTO signal (token_id, timestamp, name, source, value_number, value_string) 
VALUES 
(1, now(), 'garbagename1', 'dimo/integration/1', 0, 'garbagevalue1'),
(1, now(), 'garbagename2', 'dimo/integration/2', 0, 'garbagevalue2'),
(1, now(), 'garbagename3', 'dimo/integration/3', 0, 'garbagevalue3'),
(2, now(), 'garbagename1', 'dimo/integration/1', 0, 'garbagevalue1'),
(3, now(), 'garbagename1', 'dimo/integration/2', 0, 'garbagevalue1'), 
(4, now(), 'garbagename1', 'dimo/integration/3', 0, 'garbagevalue1'),
(5, now(), 'garbagename1', 'dimo/integration/1', 0, 'garbagevalue1'), 
(6, now(), 'garbagename1', 'dimo/integration/2', 0, 'garbagevalue1'),
(7, now(), 'garbagename1', 'dimo/integration/3', 0, 'garbagevalue1'), 
(8, now(), 'garbagename1', 'dimo/integration/1', 0, 'garbagevalue1'),
(9, now(), 'garbagename1', 'dimo/integration/2', 0, 'garbagevalue1'), 
(10, now(), 'garbagename1', 'dimo/integration/3', 0, 'garbagevalue1');
`
