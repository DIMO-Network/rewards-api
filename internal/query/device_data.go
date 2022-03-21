package query

import (
	"context"
	"time"

	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
)

type EngineType int

const (
	Gasoline EngineType = iota
	Hybrid
	Electric
)

type DeviceDataClient interface {
	GetMilesDriven(userDeviceID string, start, end time.Time) float64
	UsesFuel(userDeviceID string, start, end time.Time) bool
	UsesElectricity(userDeviceID string, start, end time.Time) bool
}

type elasticDeviceDataClient struct {
	startTime, endTime time.Time
	client             *elasticsearch.Client
}

func (c *elasticDeviceDataClient) GetMilesDriven(userDeviceID string) float64 {
	ctx := context.Background()
	query := func(order esquery.Order) *esquery.SearchRequest {
		return esquery.Search().
			Query(
				esquery.Bool().
					Must(
						esquery.Term("subject", userDeviceID),
						esquery.Exists("data.odometer"),
						esquery.Range("data.timestamp").Gte(c.startTime).Lt(c.endTime),
					),
			).
			Sort("data.timestamp", esquery.OrderAsc).
			Size(1)
	}

	query(esquery.OrderAsc).Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex("device-status-prod-*"))
	return 0
}

func (c *elasticDeviceDataClient) UsesFuel(userDeviceID string, start, end time.Time) bool {
	ctx := context.Background()
	query := func(order esquery.Order) *esquery.SearchRequest {
		return esquery.Search().
			Query(
				esquery.Bool().
					Must(
						esquery.Term("subject", userDeviceID),
						esquery.Exists("data.fuelPercentRemaining"),
						esquery.Range("data.timestamp").Gte(c.startTime).Lt(c.endTime),
					),
			).
			Size(1)
	}

	query(esquery.OrderAsc).Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex("device-status-prod-*"))
	return false
}

func (c *elasticDeviceDataClient) UsesElectricity(userDeviceID string, start, end time.Time) bool {
	ctx := context.Background()
	query := func(order esquery.Order) *esquery.SearchRequest {
		return esquery.Search().
			Query(
				esquery.Bool().
					Must(
						esquery.Term("subject", userDeviceID),
						esquery.Exists("data.fuelPercentRemaining"),
						esquery.Range("data.soc").Gte(c.startTime).Lt(c.endTime),
					),
			).
			Size(1)
	}

	query(esquery.OrderAsc).Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex("device-status-prod-*"))
	return false
}
