package identity

// TODO(elffjs): It would be nice to grab even more from Identity: the pairings,
// for example, could be pulled here.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

const jsonContentType = "application/json"

const query = `
	query GetStake($vehicleId: Int!) {
		vehicle(tokenId: $vehicleId) {
			stake {
				points
				endsAt
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
	Data struct {
		Vehicle *struct {
			Stake *struct {
				Points int       `json:"points"`
				EndsAt time.Time `json:"endsAt"`
			} `json:"stake"`
		} `json:"vehicle"`
	} `json:"data"`
	Errors []struct {
		Path       []string `json:"path"`
		Extensions struct {
			Code string `json:"code"`
		}
	} `json:"errors"`
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

	if len(resBody.Errors) != 0 {
		if len(resBody.Errors) == 1 {
			oneError := resBody.Errors[0]
			if slices.Equal([]string{"vehicle"}, oneError.Path) && oneError.Extensions.Code == "NOT_FOUND" {
				// This is actually kinda bad. Why can't we find the vehicle?
				return 0, nil
			}
		}

		return 0, fmt.Errorf("unexpected error: %v", resBody.Errors)
	}

	// Really shouldn't be possible for vehicle to be nil without an error.
	if resBody.Data.Vehicle == nil || resBody.Data.Vehicle.Stake == nil || resBody.Data.Vehicle.Stake.EndsAt.Before(time.Now()) {
		return 0, nil
	}

	return resBody.Data.Vehicle.Stake.Points, nil
}
