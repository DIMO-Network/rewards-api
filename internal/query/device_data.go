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
	GetMilesDriven(userDeviceID string, start, end time.Time) (float64, error)
	UsesFuel(userDeviceID string, start, end time.Time) (bool, error)
	UsesElectricity(userDeviceID string, start, end time.Time) (bool, error)
	ListActiveUserDeviceIDs(start, end time.Time) ([]string, error)
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

func (c *elasticDeviceDataClient) GetMilesDriven(userDeviceID string, start, end time.Time) (float64, error) {
	ctx := context.Background()
	query := func(order esquery.Order) *esquery.SearchRequest {
		return esquery.Search().
			Query(
				esquery.Bool().
					Filter(
						esquery.Term("subject", userDeviceID),
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
		return 0, err
	}
	defer resp.Body.Close()
	respB := new(odometerResp)
	if err := json.NewDecoder(resp.Body).Decode(respB); err != nil {
		fmt.Printf("%v\n", err)
		return 0, err
	}

	q = query(esquery.OrderDesc)
	resp, err = q.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0, err
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("%v\n", err)
		return 0, err
	}
	respC := new(odometerResp)
	if err := json.NewDecoder(resp.Body).Decode(respC); err != nil {
		fmt.Printf("%v\n", err)
		return 0, err
	}

	if len(respC.Hits.Hits) == 0 {
		return 0, nil
	}

	return respC.Hits.Hits[0].Source.Data.Odometer - respB.Hits.Hits[0].Source.Data.Odometer, nil
}

func (c *elasticDeviceDataClient) ListActiveUserDeviceIDs(start, end time.Time) ([]string, error) {
	ctx := context.Background()

	query := esquery.Search().
		Query(
			esquery.Bool().
				Filter(
					esquery.Exists("data.odometer"),
					esquery.Range("data.timestamp").Gte(start).Lt(end),
				),
		).
		Aggs(
			esquery.CustomAgg("active_devices",
				map[string]interface{}{
					"composite": map[string]interface{}{
						"size": 100,
						"sources": []map[string]interface{}{
							{
								"subject": map[string]interface{}{
									"terms": map[string]interface{}{
										"field": "subject",
									},
								},
							},
						},
					},
				},
			),
		).
		Size(0)

	b, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	resp, err := query.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		return nil, err
	}

	respC := new(DeviceListResp)
	if err := json.NewDecoder(resp.Body).Decode(respC); err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	out := make([]string, len(respC.Aggregations.ActiveDevices.Buckets))
	for i, dv := range respC.Aggregations.ActiveDevices.Buckets {
		out[i] = dv.Key.Subject
	}

	return out, nil
}

func (c *elasticDeviceDataClient) UsesFuel(userDeviceID string, start, end time.Time) (bool, error) {
	ctx := context.Background()

	query := esquery.Search().
		Query(
			esquery.Bool().
				Must(
					esquery.Term("subject", userDeviceID),
					esquery.Exists("data.fuelPercentRemaining"),
					esquery.Range("data.timestamp").Gte(start).Lt(end),
				),
		).
		Size(1)

	resp, err := query.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		return false, err
	}

	respC := new(odometerResp)
	if err := json.NewDecoder(resp.Body).Decode(respC); err != nil {
		fmt.Printf("%v\n", err)
		return false, err
	}

	return len(respC.Hits.Hits) > 0, nil
}

type DeviceListResp struct {
	Aggregations struct {
		ActiveDevices struct {
			Buckets []struct {
				Key struct {
					Subject string `json:"subject"`
				} `json:"key"`
			} `json:"buckets"`
		} `json:"active_devices"`
	} `json:"aggregations"`
}

func (c *elasticDeviceDataClient) UsesElectricity(userDeviceID string, start, end time.Time) (bool, error) {
	ctx := context.Background()
	query := esquery.Search().
		Query(
			esquery.Bool().
				Must(
					esquery.Term("subject", userDeviceID),
					esquery.Exists("data.soc"),
					esquery.Range("data.timestamp").Gte(start).Lt(end),
				),
		).
		Size(1)

	resp, err := query.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		return false, err
	}

	respC := new(odometerResp)
	if err := json.NewDecoder(resp.Body).Decode(respC); err != nil {
		fmt.Printf("%v\n", err)
		return false, err
	}

	return len(respC.Hits.Hits) > 0, nil
}
