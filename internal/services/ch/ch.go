package ch

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DIMO-Network/clickhouse-infra/pkg/connect"
	"github.com/DIMO-Network/model-garage/pkg/vss"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/aarondl/sqlboiler/v4/drivers"
	"github.com/aarondl/sqlboiler/v4/queries"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/aarondl/sqlboiler/v4/queries/qmhelper"
)

// Client is a ClickHouse client that interacts with the ClickHouse database.
type Client struct {
	conn clickhouse.Conn
}

// constantSignals are signals that never change. They tend to come from VIN decoding and should
// not be taken as signs of an active connection.
var constantSignals = []any{vss.FieldPowertrainTractionBatteryGrossCapacity}

var (
	dialect = drivers.Dialect{
		LQ: '`',
		RQ: '`',
	}
)

// NewClient creates a new ClickHouse client.
func NewClient(settings *config.Settings) (*Client, error) {
	conn, err := connect.GetClickhouseConn(&settings.Clickhouse)
	if err != nil {
		return nil, fmt.Errorf("failed to create ClickHouse connection: %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}
	return &Client{conn: conn}, nil
}

// DescribeActiveDevices returns a list of vehicles and associated integrations KSUIDs that have signals
// in ClickHouse during specified time period. The vehicles are sorted by token id, descending.
func (s *Client) DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*Vehicle, error) {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q,
		qm.Select("token_id", "groupUniqArray(source)"),
		qm.From("signal"),
		qmhelper.Where("timestamp", qmhelper.GTE, start),
		qmhelper.Where("timestamp", qmhelper.LT, end),
		qm.WhereNotIn("name NOT IN ?", constantSignals...),
		qm.GroupBy("token_id"),
		qm.OrderBy("token_id DESC"),
	)
	query, args := queries.BuildQuery(q)
	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clickhouse: %w", err)
	}
	defer rows.Close()

	var vehicles []*Vehicle
	for rows.Next() {
		var tokenID uint32
		var sources []string
		if err := rows.Scan(&tokenID, &sources); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		vehicles = append(vehicles, &Vehicle{
			TokenID: int64(tokenID),
			Sources: sources,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("clickhouse row error: %w", rows.Err())
	}

	return vehicles, nil
}

// GetSourcesForVehicles is like DescribeActiveDevices, restricted to a given list of vehicle
// token ids.
func (s *Client) GetSourcesForVehicles(ctx context.Context, tokenIDs []uint64, start, end time.Time) ([]*Vehicle, error) {
	if len(tokenIDs) == 0 {
		return nil, nil
	}

	// Need to convert to []any for the query mod to work.
	anyIDs := make([]any, len(tokenIDs))
	for i, x := range tokenIDs {
		anyIDs[i] = x
	}

	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q,
		qm.Select("token_id", "groupUniqArray(source)"),
		qm.From("signal"),
		qmhelper.Where("timestamp", qmhelper.GTE, start),
		qmhelper.Where("timestamp", qmhelper.LT, end),
		qm.WhereIn("token_id IN ?", anyIDs...),
		qm.WhereNotIn("name NOT IN ?", constantSignals...),
		qm.GroupBy("token_id"),
	)
	query, args := queries.BuildQuery(q)
	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clickhouse: %w", err)
	}
	defer rows.Close()

	var vehicles []*Vehicle
	for rows.Next() {
		var tokenID uint32
		var sources []string
		if err := rows.Scan(&tokenID, &sources); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		vehicles = append(vehicles, &Vehicle{
			TokenID: int64(tokenID),
			Sources: sources,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("clickhouse row error: %w", rows.Err())
	}

	return vehicles, nil
}

type Vehicle struct {
	TokenID int64
	Sources []string
}
