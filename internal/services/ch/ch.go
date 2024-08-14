package ch

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DIMO-Network/clickhouse-infra/pkg/connect"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// Client is a ClickHouse client that interacts with the ClickHouse database.
type Client struct {
	conn clickhouse.Conn
}

const (
	integPrefix = "dimo/integration/"
)

var dialect = drivers.Dialect{
	LQ: '`',
	RQ: '`',
}

// NewClient creates a new ClickHouse client.
func NewClient(settings *config.Settings) (*Client, error) {
	conn, err := connect.GetClickhouseConn(&settings.Clickhouse)
	if err != nil {
		return nil, fmt.Errorf("failed to create ClickHouse connection: %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse: %w", err)
	}
	return &Client{conn: conn}, nil
}

// DescribeActiveDevices returns list of vehicles and associated vehicle integrations that have signals in Clickhouse during specified time period
func (s *Client) DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*Vehicle, error) {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, []qm.QueryMod{
		qm.Select("token_id", "groupUniqArray(source)"),
		qm.From("signal"),
		qm.Where("timestamp >= ?", start),
		qm.And("timestamp < ?", end),
		qm.GroupBy("token_id"),
	}...)
	query, args := queries.BuildQuery(q)
	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clickhouse: %w", err)
	}
	defer rows.Close()

	var vehicles []*Vehicle
	for rows.Next() {
		var tokenID uint32
		var integrations []string
		if err := rows.Scan(&tokenID, &integrations); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		for i := range integrations {
			integrations[i] = strings.TrimPrefix(integrations[i], integPrefix)
		}

		vehicles = append(vehicles, &Vehicle{
			TokenID:      int64(tokenID),
			Integrations: integrations,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("clickhouse row error: %w", rows.Err())
	}

	return vehicles, nil
}

// GetIntegrations returns list of integrations from Clickhouse for a given vehicle
func (s *Client) GetIntegrations(ctx context.Context, tokenID uint64, start, end time.Time) ([]string, error) {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, []qm.QueryMod{
		qm.Select("groupUniqArray(source)"),
		qm.From("signal"),
		qm.Where("timestamp >= ?", start),
		qm.And("timestamp < ?", end),
		qm.And("token_id = ?", tokenID),
	}...)
	query, args := queries.BuildQuery(q)
	row := s.conn.QueryRow(ctx, query, args...)

	var vehIntegs []string
	if err := row.Scan(&vehIntegs); err != nil {
		return nil, fmt.Errorf("failed to scan clickhouse row: %w", err)
	}

	if row.Err() != nil {
		return nil, fmt.Errorf("failed to query clickhouse: %w", row.Err())
	}

	for i := range vehIntegs {
		vehIntegs[i] = strings.TrimPrefix(vehIntegs[i], integPrefix)
	}

	return vehIntegs, nil
}

type Vehicle struct {
	TokenID      int64
	Integrations []string
}
