package mobileapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
)

type Client struct {
	base *url.URL
	http *http.Client
}

var zeroAddr common.Address

func New(base *url.URL) *Client {
	return &Client{
		base: base,
		http: &http.Client{},
	}
}

type referrerResp struct {
	Address *common.Address `json:"address"`
}

var ErrNoReferrer = errors.New("no referrer found")

// GetReferrer returns the address of the user that referred the given user, using the
// Mobile API. If no referrer exists then this returns the error ErrNoReferrer.
// This uses https://api.dimo.co/api#/referral/ReferralController_getReferrer
func (c *Client) GetReferrer(ctx context.Context, addr common.Address) (common.Address, error) {
	p := c.base.JoinPath("referral", "referrer", addr.Hex())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.String(), nil)
	if err != nil {
		return zeroAddr, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return zeroAddr, err
	}

	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return zeroAddr, err
	}

	if resp.StatusCode != http.StatusOK {
		return zeroAddr, fmt.Errorf("status code %d", resp.StatusCode)
	}

	var respBody referrerResp
	err = json.Unmarshal(b, &respBody)
	if err != nil {
		return zeroAddr, err
	}

	if respBody.Address == nil {
		return zeroAddr, ErrNoReferrer
	}

	return *respBody.Address, nil
}
