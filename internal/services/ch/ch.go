package ch

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DIMO-Network/clickhouse-infra/pkg/connect"
	"github.com/DIMO-Network/rewards-api/internal/config"
)

// Client is a ClickHouse client that interacts with the ClickHouse database.
type Client struct {
	conn clickhouse.Conn
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

func (s *Client) DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*Devices, error) {
	startStr := start.Format("2006-01-02 15:04:05")
	endStr := end.Format("2006-01-02 15:04:05")
	query := fmt.Sprintf("SELECT token_id, groupUniqArray(source) FROM signal WHERE timestamp >= '%s' AND timestamp < '%s' GROUP BY token_id;", startStr, endStr)

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query clickhouse: %w", err)
	}
	defer rows.Close()

	var devices []*Devices
	for rows.Next() {
		var tokenID uint32
		var integrations []string
		if err := rows.Scan(&tokenID, &integrations); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}
		devices = append(devices, &Devices{
			VehicleTokenID: int64(tokenID),
			Integrations:   integrations,
		})
	}

	return devices, nil
}

type Devices struct {
	VehicleTokenID int64
	Integrations   []string
}
