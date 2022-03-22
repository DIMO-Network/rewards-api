package query

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
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
	client *elasticsearch.Client
	index  string
}

func NewDeviceDataClient(settings *config.Settings) DeviceDataClient {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses:            []string{settings.ElasticSearchAnalyticsHost},
			Username:             settings.ElasticSearchAnalyticsUsername,
			Password:             settings.ElasticSearchAnalyticsPassword,
			EnableRetryOnTimeout: true,
			MaxRetries:           5,
		},
	)
	if err != nil {
		panic(err)
	}
	return &elasticDeviceDataClient{
		client: client,
		index:  settings.DeviceDataIndexName,
	}
}

type odometerResp struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Data struct {
					Odometer float64 `json:"odometer"`
				} `json:"data"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (c *elasticDeviceDataClient) GetMilesDriven(userDeviceID string, start, end time.Time) float64 {
	ctx := context.Background()
	fmt.Println(start.UTC(), end.UTC())
	query := func(order esquery.Order) *esquery.SearchRequest {
		return esquery.Search().
			Query(
				esquery.Bool().
					Filter(
						esquery.Term("subject", "AX"),
						esquery.Exists("data.odometer"),
						esquery.Range("data.timestamp").Gte(start).Lt(end),
					),
			).
			Sort("data.timestamp", order).
			Size(1)
	}

	q := query(esquery.OrderAsc)
	resp, err := q.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}
	defer resp.Body.Close()
	respB := new(odometerResp)
	if err := json.NewDecoder(resp.Body).Decode(respB); err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}

	q = query(esquery.OrderDesc)
	resp, err = q.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}
	respC := new(odometerResp)
	if err := json.NewDecoder(resp.Body).Decode(respC); err != nil {
		fmt.Printf("%v\n", err)
		return 0
	}

	return respC.Hits.Hits[0].Source.Data.Odometer - respB.Hits.Hits[0].Source.Data.Odometer
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
						esquery.Range("data.timestamp").Gte(start).Lt(end),
					),
			).
			Size(1)
	}

	query(esquery.OrderAsc).Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
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
						esquery.Range("data.soc").Gte(start).Lt(end),
					),
			).
			Size(1)
	}

	query(esquery.OrderAsc).Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	return false
}
