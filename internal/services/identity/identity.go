package identity

// TODO(elffjs): It would be nice to grab even more from Identity: the pairings,
// for example, could be pulled here.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

type VehicleDescr struct {
	Owner             common.Address `json:"owner"`
	AftermarketDevice *struct {
		TokenID      int            `json:"tokenId"`
		Owner        common.Address `json:"owner"`
		Beneficiary  common.Address `json:"beneficiary"`
		Manufacturer struct {
			TokenID int `json:"tokenId"`
		} `json:"manufacturer"`
	} `json:"aftermarketDevice"`
	SyntheticDevice *struct {
		TokenID    int `json:"tokenId"`
		Connection struct {
			Address common.Address `json:"address"`
		} `json:"connection"`
	} `json:"syntheticDevice"`
	Stake *struct {
		Points int       `json:"points"`
		EndsAt time.Time `json:"endsAt"`
	} `json:"stake"`
}

type resp struct {
	Data struct {
		Vehicle *VehicleDescr `json:"vehicle"`
	} `json:"data"`
	Errors []struct {
		Path       []string `json:"path"`
		Extensions struct {
			Code string `json:"code"`
		}
	} `json:"errors"`
}

func (c *Client) DescribeVehicle(vehicleID int) (*VehicleDescr, error) {
	p := payload{
		Query: query,
		Variables: map[string]any{
			"vehicleId": vehicleID,
		},
	}

	reqBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal request body: %w", err)
	}

	res, err := c.Client.Post(c.QueryURL, jsonContentType, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %w", err)
	}

	var resBody resp

	err = json.Unmarshal(resBytes, &resBody)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse response body: %w", err)
	}

	if len(resBody.Errors) != 0 {
		if len(resBody.Errors) == 1 {
			oneError := resBody.Errors[0]
			if slices.Equal([]string{"vehicle"}, oneError.Path) && oneError.Extensions.Code == "NOT_FOUND" {
				// This is actually kinda bad. Why can't we find the vehicle?
				return nil, ErrNotFound
			}
		}

		return nil, fmt.Errorf("unexpected error: %v", resBody.Errors)
	}

	// Really shouldn't be possible for vehicle to be nil without an error.
	if resBody.Data.Vehicle == nil {
		return nil, fmt.Errorf("no error, but vehicle response is null")
	}

	return resBody.Data.Vehicle, nil
}

var ErrNotFound = errors.New("no vehicle with that token id found")
