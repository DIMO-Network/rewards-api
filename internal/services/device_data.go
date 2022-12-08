package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
)

type DeviceDataClient interface {
	GetLastActivity(userDeviceID string) (lastActivity time.Time, seen bool, err error)
	DescribeActiveDevices(start, end time.Time) ([]*DeviceData, error)
	GetIntegrations(userDeviceID string, start, end time.Time) (ints []string, err error)
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

type DeviceData struct {
	ID           string
	Integrations []string
	Serials      []string
}

// activityFieldExists is a list of Elastic exists queries for all of the non-constant fields in
// the device status update.
var activityFieldExists = []esquery.Mappable{
	esquery.Exists("data.odometer"),
	esquery.Exists("data.soc"),
	esquery.Exists("data.latitude"),
	esquery.Exists("data.oil"),
	esquery.Exists("data.fuelPercentRemaining"),
	esquery.Exists("data.tires.frontLeft"),
	esquery.Exists("data.charging"),
}

func (c *elasticDeviceDataClient) GetLastActivity(userDeviceID string) (lastActivity time.Time, seen bool, err error) {
	ctx := context.Background()

	query := esquery.Search().
		Query(
			esquery.Bool().
				Filter(esquery.Term("subject", userDeviceID)).
				Should(activityFieldExists...).
				MinimumShouldMatch(1),
		).
		SourceIncludes("data.timestamp").
		Sort("data.timestamp", esquery.OrderDesc).
		Size(1)

	res, err := query.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		return time.Time{}, false, err
	}
	defer res.Body.Close()

	respb := new(ActivityResp)
	if err := json.NewDecoder(res.Body).Decode(respb); err != nil {
		return time.Time{}, false, err
	}

	if len(respb.Hits.Hits) == 0 {
		return time.Time{}, false, nil
	}

	return respb.Hits.Hits[0].Source.Data.Timestamp, true, nil
}

type ActivityResp struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Data struct {
					Timestamp time.Time `json:"timestamp"`
				} `json:"data"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (c *elasticDeviceDataClient) GetIntegrations(userDeviceID string, start, end time.Time) (ints []string, err error) {
	ctx := context.Background()

	query := esquery.Search().
		Query(
			esquery.Bool().
				Filter(
					esquery.Term("subject", userDeviceID),
					esquery.Range("data.timestamp").Gte(start).Lt(end),
				).
				Should(activityFieldExists...).
				MinimumShouldMatch(1),
		).
		Aggs(
			esquery.TermsAgg("integrations", "source"),
		).
		Size(0)

	qb, _ := json.Marshal(query)
	log.Print(string(qb))

	res, err := query.Run(c.client, c.client.Search.WithContext(ctx), c.client.Search.WithIndex(c.index))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	log.Print("code", res.StatusCode)

	if code := res.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("status code %d", code)
	}

	rb, _ := io.ReadAll(res.Body)
	log.Print(string(rb))

	respb := new(DeviceIntegrationsResp)
	if err := json.Unmarshal(rb, respb); err != nil {
		return nil, err
	}

	integrations := []string{}
	for _, bucket := range respb.Aggregations.Integrations.Buckets {
		key := bucket.Key
		if strings.HasPrefix(key, integPrefix) {
			integrations = append(integrations, strings.TrimPrefix(key, integPrefix))
		}
	}

	return integrations, nil
}

type DeviceIntegrationsResp struct {
	Aggregations struct {
		Integrations struct {
			Buckets []struct {
				Key string `json:"key"`
			} `json:"buckets"`
		} `json:"integrations"`
		UnitIDs struct {
			Buckets []struct {
				Key string `json:"key"`
			} `json:"buckets"`
		} `json:"unit_ids"`
	} `json:"aggregations"`
}

func (c *elasticDeviceDataClient) DescribeActiveDevices(start, end time.Time) ([]*DeviceData, error) {
	ctx := context.Background()

	out := make([]*DeviceData, 0)
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
					Filter(esquery.Range("data.timestamp").Gte(start).Lt(end)).
					Should(activityFieldExists...).
					MinimumShouldMatch(1),
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
			return nil, err
		}

		for _, dv := range respBody.Aggregations.ActiveDevices.Buckets {
			deviceID := dv.Key.Subject
			integrations := make([]string, 0)
			for _, bucket := range dv.Integrations.Buckets {
				key := bucket.Key
				if strings.HasPrefix(key, integPrefix) {
					integrations = append(integrations, strings.TrimPrefix(key, integPrefix))
				}
			}
			serials := make([]string, 0)
			for _, bucket := range dv.UnitIDs.Buckets {
				key := bucket.Key
				serials = append(serials, key)
			}
			out = append(out, &DeviceData{
				ID:           deviceID,
				Integrations: integrations,
				Serials:      serials,
			})
		}

		if respBody.Aggregations.ActiveDevices.AfterKey.Subject != nil {
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
				UnitIDs struct {
					Buckets []struct {
						Key string `json:"key"`
					} `json:"buckets"`
				} `json:"unit_ids"`
			} `json:"buckets"`
		} `json:"active_devices"`
	} `json:"aggregations"`
}
