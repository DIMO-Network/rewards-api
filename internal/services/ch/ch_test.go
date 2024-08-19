package ch

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"math/rand/v2"

	chconfig "github.com/DIMO-Network/clickhouse-infra/pkg/connect/config"
	"github.com/DIMO-Network/clickhouse-infra/pkg/container"
	dimo_ch "github.com/DIMO-Network/clickhouse-infra/pkg/container"
	"github.com/DIMO-Network/model-garage/pkg/migrations"
	"github.com/stretchr/testify/suite"
)

type CHTestSuite struct {
	suite.Suite
	chClient       *Client
	container      *container.Container
	tokenSourceMap map[int64][]string
}

var sources = []string{"2ULfuC8U9dOqRshZBAi0lMM1Rrx", "27qftVRWQYpVDcO5DltO5Ojbjxk", "22N2xaPOq2WW2gAHBHd0Ikn4Zob"}

const dummyTokens int64 = 10

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

	batch, err := conn.PrepareBatch(ctx, fmt.Sprintf("INSERT INTO signal"))
	c.Require().NoError(err, "Failed to prepare batch")

	c.tokenSourceMap = make(map[int64][]string)
	dummyData := generateRandomData(dummyTokens)
	for _, d := range dummyData {
		err := batch.AppendStruct(&d)
		c.Require().NoError(err, "Failed to append struct")
		if _, ok := c.tokenSourceMap[d.TokenID]; !ok {
			c.tokenSourceMap[d.TokenID] = []string{}
		}
		c.tokenSourceMap[d.TokenID] = append(c.tokenSourceMap[d.TokenID], strings.TrimPrefix(d.Source, integPrefix))
	}

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
	start := time.Now().AddDate(0, 0, -6)
	end := time.Now().AddDate(0, 0, 1)
	resp, err := c.chClient.DescribeActiveDevices(ctx, start, end)
	c.Require().NoError(err)
	c.Require().Equal(len(c.tokenSourceMap), len(resp))
	for _, r := range resp {
		c.Require().ElementsMatch(c.tokenSourceMap[r.TokenID], r.Integrations)
	}
}

func (c *CHTestSuite) Test_GetIntegrations() {
	ctx := context.Background()
	start := time.Now().AddDate(0, 0, -6)
	end := time.Now().AddDate(0, 0, 1)

	for k, s := range c.tokenSourceMap {
		resp, err := c.chClient.GetIntegrations(ctx, uint64(k), start, end)
		c.Require().NoError(err)
		c.Require().ElementsMatch(s, resp)
	}
}

type data struct {
	TokenID   int64     `ch:"token_id"`
	Timestamp time.Time `ch:"timestamp"`
	Name      string    `ch:"name"`
	Source    string    `ch:"source"`
	ValueN    float64   `ch:"value_number"`
	ValueS    string    `ch:"value_string"`
}

func generateRandomData(all int64) []data {
	d := []data{}
	for tokenID := range all {
		for n, s := range sources[:rand.IntN(len(sources))] {
			d = append(d, data{
				TokenID:   tokenID,
				Timestamp: time.Now(),
				Name:      fmt.Sprintf("name%d", n),
				Source:    fmt.Sprintf("dimo/integration/%s", s),
				ValueS:    fmt.Sprintf("value%d", n),
			})
		}
	}

	return d
}
