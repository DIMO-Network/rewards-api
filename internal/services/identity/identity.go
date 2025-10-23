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

const vehicleFragment = `
fragment VehicleInfo on Vehicle {
	tokenId
	owner
	aftermarketDevice {
		tokenId
		beneficiary
		manufacturer {
			tokenId
		}
	}
	syntheticDevice {
		tokenId
		connection {
			address
		}
	}
	stake {
		points
		endsAt
	}
}
`

const query = vehicleFragment + `
	query DescribeVehicle($vehicleId: Int!) {
		vehicle(tokenId: $vehicleId) {
			...VehicleInfo
		}
	}`

const ownerQuery = vehicleFragment + `
	query OwnedVehicles($owner: Address!, $after: String) {
		vehicles(first: 100, after: $after, filterBy: {owner: $owner}) {
			pageInfo {
				hasNextPage
				endCursor
			}
			nodes {
				...VehicleInfo
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

type VehicleDescription struct {
	TokenID           int            `json:"tokenId"`
	Owner             common.Address `json:"owner"`
	AftermarketDevice *struct {
		TokenID      int            `json:"tokenId"`
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
		Vehicle *VehicleDescription `json:"vehicle"`
	} `json:"data"`
	Errors []struct {
		Path       []string `json:"path"`
		Extensions struct {
			Code string `json:"code"`
		}
	} `json:"errors"`
}

type listResp struct {
	Data struct {
		Vehicles struct {
			PageInfo struct {
				HasNextPage bool    `json:"hasNextPage"`
				EndCursor   *string `json:"endCursor"`
			} `json:"pageInfo"`
			Nodes []*VehicleDescription `json:"nodes"`
		} `json:"vehicles"`
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
			"vehicleId": vehicleID,
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

func (c *Client) DescribeVehicle(vehicleID uint64) (*VehicleDescription, error) {
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
				return nil, err
			}
		}

		return nil, fmt.Errorf("unexpected error: %v", resBody.Errors)
	}

	// Really shouldn't be possible for vehicle to be nil without an error.
	if resBody.Data.Vehicle == nil {
		return nil, errors.New("no error, but vehicle is empty")
	}

	return resBody.Data.Vehicle, nil
}

var ErrNotFound = errors.New("no vehicle with that token id found")

func (c *Client) GetVehicles(owner common.Address) ([]*VehicleDescription, error) {
	var after *string

	var out []*VehicleDescription

	for {
		p := payload{
			Query: ownerQuery,
			Variables: map[string]any{
				"owner": owner,
				"after": after,
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

		var resBody listResp

		err = json.Unmarshal(resBytes, &resBody)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse response body: %w", err)
		}

		if len(resBody.Errors) != 0 {
			return nil, fmt.Errorf("unexpected errors: %v", resBody.Errors)
		}

		out = append(out, resBody.Data.Vehicles.Nodes...)

		if resBody.Data.Vehicles.PageInfo.HasNextPage {
			after = resBody.Data.Vehicles.PageInfo.EndCursor
		} else {
			break
		}
	}

	return out, nil
}
