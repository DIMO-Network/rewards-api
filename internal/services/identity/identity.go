package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const jsonContentType = "application/json"

const query = `{
	query GetStake($vehicleId: Int!) {
		vehicle(by: {tokenId: $vehicleId}) {
			stake {
				points
				endsAt
			}
		}
	}
}`

type Client struct {
	QueryURL string
	Client   *http.Client
}

type payload struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type resp struct {
	Vehicle *struct {
		Stake *struct {
			Points int       `json:"points"`
			EndsAt time.Time `json:"endsAt"`
		} `json:"stake"`
	} `json:"vehicle"`
}

// GetVehicleStakePoints returns the number of points that should be added to the vehicle's weekly
// total because of $DIMO staking. This number will be 0 if the vehicle has no attached stake, or
// if the attached stake has ended.
func (c *Client) GetVehicleStakePoints(vehicleID uint64) (int, error) {
	p := payload{
		Query: query,
		Variables: map[string]any{
			"tokenId": vehicleID,
		},
	}

	reqBytes, err := json.Marshal(p)
	if err != nil {
		return 0, fmt.Errorf("couldn't marshal request body: %w", err)
	}

	res, err := c.Client.Post(c.QueryURL, jsonContentType, bytes.NewBuffer(reqBytes))
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code %d", res.StatusCode)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("couldn't read response body: %w", err)
	}

	var resBody resp

	err = json.Unmarshal(resBytes, &resBody)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse response body: %w", err)
	}

	// TODO(elffjs): The error handling here is too loose: if this failed because of, e.g., a
	// database issue then we want to bail.
	if resBody.Vehicle == nil || resBody.Vehicle.Stake == nil {
		return 0, nil
	}

	return resBody.Vehicle.Stake.Points, nil
}
