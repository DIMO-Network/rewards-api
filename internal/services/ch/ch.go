package ch

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DIMO-Network/rewards-api/internal/config"
)

const (
	defaultMaxExecutionTime                    = 5
	defaultTimeoutBeforeCheckingExecutionSpeed = 2
	// TimeoutErrCode is the error code returned by ClickHouse when a query is interrupted due to exceeding the max_execution_time.
	TimeoutErrCode = int32(159)
)

// Client is a ClickHouse client that interacts with the ClickHouse database.
type Client struct {
	conn clickhouse.Conn
}

// NewClient creates a new ClickHouse client.
func NewClient(settings *config.Settings) (*Client, error) {
	maxExecutionTime, err := getMaxExecutionTime(settings.MaxRequestDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to get max execution time: %w", err)
	}
	addr := fmt.Sprintf("%s:%d", settings.Clickhouse.Host, settings.Clickhouse.Port)
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Username: settings.Clickhouse.User,
			Password: settings.Clickhouse.Password,
			Database: settings.Clickhouse.Database,
		},
		TLS: &tls.Config{
			RootCAs: settings.Clickhouse.RootCAs,
		},
		Settings: map[string]any{
			// ClickHouse will interrupt a query if the projected execution time exceeds the specified max_execution_time.
			// The estimated execution time is calculated after `timeout_before_checking_execution_speed`
			// More info: https://clickhouse.com/docs/en/operations/settings/query-complexity#max-execution-time
			"max_execution_time":                      maxExecutionTime,
			"timeout_before_checking_execution_speed": defaultTimeoutBeforeCheckingExecutionSpeed,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open clickhouse connection: %w", err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping clickhouse: %w", err)
	}
	return &Client{conn: conn}, nil
}

func getMaxExecutionTime(maxRequestDuration string) (int, error) {
	if maxRequestDuration == "" {
		return defaultMaxExecutionTime, nil
	}
	maxExecutionTime, err := time.ParseDuration(maxRequestDuration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse max request duration: %w", err)
	}
	return int(maxExecutionTime.Seconds()), nil
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
