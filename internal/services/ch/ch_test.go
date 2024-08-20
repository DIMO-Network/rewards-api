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
	chClient  *Client
	container *container.Container
	tokenMap  map[int]tknDtValues
}

var sources = []string{"2ULfuC8U9dOqRshZBAi0lMM1Rrx", "27qftVRWQYpVDcO5DltO5Ojbjxk", "22N2xaPOq2WW2gAHBHd0Ikn4Zob"}
var timeRanges = map[int]tknDtValues{
	1: {start: time.Now().AddDate(0, 0, -6), end: time.Now().AddDate(0, 0, 1)},
	2: {start: time.Now().AddDate(0, 0, -14), end: time.Now().AddDate(0, 0, -7)},
}

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

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO signal")
	c.Require().NoError(err, "Failed to prepare batch")

	c.tokenMap = map[int]tknDtValues{
		1: {start: time.Now().AddDate(0, 0, -6), end: time.Now().AddDate(0, 0, 1), values: map[int64][]string{}},
		2: {start: time.Now().AddDate(0, 0, -14), end: time.Now().AddDate(0, 0, -7), values: map[int64][]string{}},
	}

	for n, daterange := range c.tokenMap {
		dummyData := generateRandomData(dummyTokens, daterange)
		for _, d := range dummyData {
			err := batch.AppendStruct(&d)
			c.Require().NoError(err, "Failed to append struct")
			if _, ok := c.tokenMap[n].values[d.TokenID]; !ok {
				c.tokenMap[n].values[d.TokenID] = []string{}
			}
			c.tokenMap[n].values[d.TokenID] = append(c.tokenMap[n].values[d.TokenID], strings.TrimPrefix(d.Source, integPrefix))
		}
	}

	// for k, v := range c.tokenMap {
	// 	fmt.Println(
	// 		v.start.Format("2006-01-02"), v.end.Format("2006-01-02"),
	// 	)
	// }

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
	for n, v := range c.tokenMap {
		resp, err := c.chClient.DescribeActiveDevices(ctx, v.start, v.end)
		c.Require().NoError(err)
		c.Require().Equal(len(c.tokenMap[n].values), len(resp))
		for _, r := range resp {
			c.Require().ElementsMatch(c.tokenMap[n].values[r.TokenID], r.Integrations)
		}
	}
}

func (c *CHTestSuite) Test_GetIntegrations() {
	ctx := context.Background()
	for _, v := range c.tokenMap {
		for tkn, sources := range v.values {
			resp, err := c.chClient.GetIntegrations(ctx, uint64(tkn), v.start, v.end)
			c.Require().NoError(err)
			c.Require().ElementsMatch(sources, resp)
		}
	}
}

func generateRandomData(allTokens int64, dateRangeValues tknDtValues) []data {
	d := []data{}

	for tokenID := range allTokens {
		for n, s := range sources[:rand.IntN(len(sources))] {
			d = append(d, data{
				TokenID:   tokenID,
				Timestamp: dateRangeValues.randomDate(),
				Name:      fmt.Sprintf("name%d", n),
				Source:    fmt.Sprintf("dimo/integration/%s", s),
				ValueS:    fmt.Sprintf("value%d", n),
			})
		}
	}

	return d
}

type data struct {
	TokenID   int64     `ch:"token_id"`
	Timestamp time.Time `ch:"timestamp"`
	Name      string    `ch:"name"`
	Source    string    `ch:"source"`
	ValueN    float64   `ch:"value_number"`
	ValueS    string    `ch:"value_string"`
}

type tknDtValues struct {
	start  time.Time
	end    time.Time
	values map[int64][]string
}

func (d *tknDtValues) randomDate() time.Time {
	delta := d.end.Unix() - d.start.Unix()

	sec := rand.Int64N(delta) + d.start.Unix()
	return time.Unix(sec, 0)
}
