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
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
)

// Client is a ClickHouse client that interacts with the ClickHouse database.
type Client struct {
	conn clickhouse.Conn
}

const integPrefix = "dimo/integration/"

var (
	dialect = drivers.Dialect{
		LQ: '`',
		RQ: '`',
	}
	connectionIDToIntegrationID = map[string]string{
		"0x5a87788D90f0ded17A35E4BDaCb47f1993021630": "2lcaMFuCO0HJIUfdq8o780Kx5n3", // ruptela
		// "0x762Fd53c4973075a6fC0d17237BC65E183299980": "27qftVRWQYpVDcO5DltO5Ojbjxk", // autopi
		// "0x354574EC2cC27A29410df751e43723B1bC362Ce4": "26A5Dk3vvvQutjSyF0Jka2DP5lg", // tesla
		// "0xd3dA3efd882BA2357709328750Cf409A9131b820": "2ULfuC8U9dOqRshZBAi0lMM1Rrx", // macaron
		// "0x25229D85599603351Ac4f7606DB92Cf85D7F6A1F": "22N2xaPOq2WW2gAHBHd0Ikn4Zob", // smartcar
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
		return nil, fmt.Errorf("failed to ping clickhouse: %w", err)
	}
	return &Client{conn: conn}, nil
}

// convertSourceToIntegrationID converts a "source" in ClickHouse to an old-style integration id.
// For AutoPi, Macaron, Smartcar, Tesla this is just stripping off a prefix; but Ruptela
// in ClickHouse goes by a "connection id", an address.
func (s *Client) convertSourceToIntegrationID(source string) string {
	if integrationID, ok := connectionIDToIntegrationID[source]; ok {
		return integrationID
	}
	return strings.TrimPrefix(source, integPrefix)
}

// DescribeActiveDevices returns list of vehicles and associated vehicle integrations that have signals in Clickhouse during specified time period
func (s *Client) DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*Vehicle, error) {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q,
		qm.Select("token_id", "groupUniqArray(source)"),
		qm.From("signal"),
		qmhelper.Where("timestamp", qmhelper.GTE, start),
		qmhelper.Where("timestamp", qmhelper.LT, end),
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
		var integrations []string
		if err := rows.Scan(&tokenID, &integrations); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		for i := range integrations {
			integrations[i] = s.convertSourceToIntegrationID(integrations[i])
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
		var integrations []string
		if err := rows.Scan(&tokenID, &integrations); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		for i := range integrations {
			integrations[i] = s.convertSourceToIntegrationID(integrations[i])
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

type Vehicle struct {
	TokenID      int64
	Integrations []string
}
