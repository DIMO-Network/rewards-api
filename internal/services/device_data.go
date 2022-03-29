package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
)

type DeviceDataClient interface {
	DescribeActiveDevices(start, end time.Time) ([]*Return, error)
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

type jsonMap map[string]interface{}

type Return struct {
	ID           string
	Driven       float64
	Integrations []string
}

func (c *elasticDeviceDataClient) DescribeActiveDevices(start, end time.Time) ([]*Return, error) {
	ctx := context.Background()

	out := make([]*Return, 0)
	var afterKey *string

	for {
		compMap := jsonMap{
			"size": 100,
			"sources": []jsonMap{
				{
					"subject": jsonMap{
						"terms": jsonMap{
							"field": "subject",
						},
					},
				},
			},
		}

		if afterKey != nil {
			compMap["after"] = jsonMap{
				"subject": *afterKey,
			}
		}

		query := esquery.Search().
			Query(
				esquery.Bool().
					Filter(
						esquery.Exists("data.odometer"),
						esquery.Range("data.timestamp").Gte(start).Lt(end),
					),
			).
			Aggs(
				esquery.CustomAgg("active_devices", jsonMap{
					"composite": compMap,
					"aggs": jsonMap{
						"integrations": jsonMap{
							"terms": jsonMap{
								"field": "source",
							},
						},
						"with_odometer": jsonMap{
							"filter": jsonMap{
								"exists": jsonMap{
									"field": "data.odometer",
								},
							},
							"aggs": jsonMap{
								"odometer_high": jsonMap{
									"top_hits": jsonMap{
										"sort": []jsonMap{
											jsonMap{
												"data.timestamp": jsonMap{
													"order": "desc",
												},
											},
										},
										"_source": jsonMap{
											"includes": []string{"data.odometer"},
										},
										"size": 1,
									},
								},
								"odometer_low": jsonMap{
									"top_hits": jsonMap{
										"sort": []jsonMap{
											jsonMap{
												"data.timestamp": jsonMap{
													"order": "asc",
												},
											},
										},
										"_source": jsonMap{
											"includes": []string{"data.odometer"},
										},
										"size": 1,
									},
								},
							},
						},
					},
				}),
			).
			Size(0)

		resp, err := query.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if code := resp.StatusCode; code != http.StatusOK {
			return nil, fmt.Errorf("status code %d", code)
		}

		respBody := new(DeviceListResp)
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			fmt.Printf("%v\n", err)
			return nil, err
		}

		for _, dv := range respBody.Aggregations.ActiveDevices.Buckets {
			deviceID := dv.Key.Subject
			var drivenKM float64
			highHits := dv.WithOdometer.OdometerHigh.Hits.Hits
			lowHits := dv.WithOdometer.OdometerLow.Hits.Hits
			if len(highHits) > 0 && len(lowHits) > 0 {
				drivenKM = highHits[0].Source.Data.Odometer - lowHits[0].Source.Data.Odometer
			}
			integrations := make([]string, 0)
			for _, bucket := range dv.Integrations.Buckets {
				key := bucket.Key
				if strings.HasPrefix(key, integPrefix) {
					integrations = append(integrations, strings.TrimPrefix(key, integPrefix))
				}
			}
			out = append(out, &Return{
				ID:           deviceID,
				Driven:       drivenKM,
				Integrations: integrations,
			})
		}

		if respBody.Aggregations.ActiveDevices.AfterKey.Subject != nil {
			fmt.Println("AFTER KEY", *respBody.Aggregations.ActiveDevices.AfterKey.Subject)
			afterKey = respBody.Aggregations.ActiveDevices.AfterKey.Subject
		} else {
			break
		}
	}

	return out, nil
}

const integPrefix = "dimo/integration/"

type DeviceListResp struct {
	Aggregations struct {
		ActiveDevices struct {
			AfterKey struct {
				Subject *string `json:"subject"`
			} `json:"after_key"`
			Buckets []struct {
				Key struct {
					Subject string `json:"subject"`
				} `json:"key"`
				Integrations struct {
					Buckets []struct {
						Key string `json:"key"`
					} `json:"buckets"`
				} `json:"integrations"`
				WithOdometer struct {
					OdometerHigh struct {
						Hits struct {
							Hits []struct {
								Source struct {
									Data struct {
										Odometer float64 `json:"odometer"`
									} `json:"data"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"odometer_high"`
					OdometerLow struct {
						Hits struct {
							Hits []struct {
								Source struct {
									Data struct {
										Odometer float64 `json:"odometer"`
									} `json:"data"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"odometer_low"`
				} `json:"with_odometer"`
			} `json:"buckets"`
		} `json:"active_devices"`
	} `json:"aggregations"`
}
