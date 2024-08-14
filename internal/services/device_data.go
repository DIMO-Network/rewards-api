package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
)

const integPrefix = "dimo/integration/"

type DeviceDataClient interface {
	DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*ch.Vehicle, error)
	GetIntegrations(userDeviceID string, start, end time.Time) (ints []string, err error)
}

type deviceDataClient struct {
	chClient *ch.Client
	esClient *elasticsearch.Client
	index    string
}

type DeviceIntegrationsResp struct {
	Aggregations struct {
		Integrations struct {
			Buckets []struct {
				Key string `json:"key"`
			} `json:"buckets"`
		} `json:"integrations"`
	} `json:"aggregations"`
}

func NewDeviceDataClient(settings *config.Settings, ch *ch.Client) DeviceDataClient {
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
	return &deviceDataClient{
		esClient: client,
		index:    settings.DeviceDataIndexName,
		chClient: ch,
	}
}

// activityFieldExists is a list of Elastic exists queries for all of the non-constant fields in
// the device status update.
var activityFieldExists = []esquery.Mappable{
	esquery.Exists("data.odometer"),
	esquery.Exists("data.soc"),
	esquery.Exists("data.latitude"),
	esquery.Exists("data.oil"),
	esquery.Exists("data.range"),
	esquery.Exists("data.speed"),
	esquery.Exists("data.batteryVoltage"),
	esquery.Exists("data.coolantTemp"),
	esquery.Exists("data.engineLoad"),
	esquery.Exists("data.barometricPressure"),
	esquery.Exists("data.intakeTemp"),
	esquery.Exists("data.runTime"),
	esquery.Exists("data.ambientTemp"),
	esquery.Exists("data.fuelPercentRemaining"),
	esquery.Exists("data.tires.frontLeft"),
	esquery.Exists("data.charging"),
	esquery.Exists("data.engineSpeed"),
	esquery.Exists("data.throttlePosition"),
}

func (c *deviceDataClient) GetIntegrations(userDeviceID string, start, end time.Time) (ints []string, err error) {
	ctx := context.Background()

	query := esquery.Search().
		Query(
			esquery.Bool().
				Filter(
					esquery.Term("subject", userDeviceID),
					esquery.Range("time").Gte(start).Lt(end),
				).
				Should(activityFieldExists...).
				MinimumShouldMatch(1),
		).
		Aggs(
			esquery.TermsAgg("integrations", "source"),
		).
		Size(0)

	res, err := query.Run(c.esClient, c.esClient.Search.WithContext(ctx), c.esClient.Search.WithIndex(c.index))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if code := res.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("status code %d", code)
	}

	var respb DeviceIntegrationsResp
	if err := json.NewDecoder(res.Body).Decode(&respb); err != nil {
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

func (c *deviceDataClient) DescribeActiveDevices(ctx context.Context, start, end time.Time) ([]*ch.Vehicle, error) {
	devices, err := c.chClient.DescribeActiveDevices(ctx, start, end)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
