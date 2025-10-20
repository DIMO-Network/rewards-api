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

// const integPrefix = "dimo/integration/"

// constantSignals are signals that never change. They tend to come from VIN decoding and should
// not be taken as signs of an active connection.
var constantSignals = []any{vss.FieldPowertrainTractionBatteryGrossCapacity}

var (
	dialect = drivers.Dialect{
		LQ: '`',
		RQ: '`',
	}
	// connectionIDToIntegrationID = map[string]string{
	// 	"0xF26421509Efe92861a587482100c6d728aBf1CD0": "2lcaMFuCO0HJIUfdq8o780Kx5n3", // ruptela
	// 	"0x5e31bBc786D7bEd95216383787deA1ab0f1c1897": "27qftVRWQYpVDcO5DltO5Ojbjxk", // autopi
	// 	"0xc4035Fecb1cc906130423EF05f9C20977F643722": "26A5Dk3vvvQutjSyF0Jka2DP5lg", // tesla
	// 	"0x4c674ddE8189aEF6e3b58F5a36d7438b2b1f6Bc2": "2ULfuC8U9dOqRshZBAi0lMM1Rrx", // macaron
	// 	"0xcd445F4c6bDAD32b68a2939b912150Fe3C88803E": "22N2xaPOq2WW2gAHBHd0Ikn4Zob", // smartcar
	// 	"0x55BF1c27d468314Ea119CF74979E2b59F962295c": "2szgr5WqMQtK2ZFM8F8qW8WUfJa", // compass
	// 	"0x5879B43D88Fa93CE8072d6612cBc8dE93E98CE5d": "2wUrxPhSUPUbiHpQjPn0UrkjuvY", // motorq
	// }
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

// // convertSourceToIntegrationID converts a "source" in ClickHouse to an old-style integration KSUID.
// // "Modern" connections will always use the the connection's Ethereum address, but older records and
// // some legacy connections like Smartcar will have sources that are a "dimo/integration/" prefix in
// // front of the KSUID.
// func (s *Client) convertSourceToIntegrationID(source string) string {
// 	if integrationID, ok := connectionIDToIntegrationID[source]; ok {
// 		return integrationID
// 	}
// 	return strings.TrimPrefix(source, integPrefix)
// }

// DescribeActiveDevices returns a list of vehicles and associated integrations KSUIDs that have signals
// in ClickHouse during specified time period.
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
			TokenID: int(tokenID),
			Sources: sources,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("clickhouse row error: %w", rows.Err())
	}

	return vehicles, nil
}

// GetIntegrationsForVehicles is like DescribeActiveDevices, restricted to a given list of vehicle
// token ids.
func (s *Client) GetIntegrationsForVehicles(ctx context.Context, tokenIDs []uint64, start, end time.Time) ([]*Vehicle, error) {
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
			TokenID: int(tokenID),
			Sources: sources,
		})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("clickhouse row error: %w", rows.Err())
	}

	return vehicles, nil
}

type Vehicle struct {
	TokenID int
	// Sources are the connection addresses. 0x-prefixed and checksummed.
	Sources []string
}
