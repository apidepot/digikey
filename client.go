// Copyright (c) 2025 The digikey developers. All rights reserved.
// Project site: https://github.com/apidepot/digikey
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package digikey

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	apiURL          = "https://api.digikey.com/v1/"
	sandboxURL      = "https://sandbox-api.digikey.com/v1/"
	accessTokenURL  = "https://api.digikey.com/v1/oauth2/token"
	sandboxTokenURL = "https://sandbox-api.digikey.com/v1/oauth2/token"
	grantType       = "client_credentials"
)

// Client models a client to consume the DigiKey API.
type Client struct {
	baseURL        string
	accessTokenURL string
	id             string
	secret         string
	accessToken    string
	tokenType      string
	tokenExpiresAt time.Time
	httpClient     *http.Client
	rateLimiter    *rate.Limiter
	mu             sync.RWMutex
}

// Error represents an IEX API error
type Error struct {
	StatusCode int
	Message    string
}

// ClientOption applies an option to the client.
type ClientOption func(*Client)

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("%d %s: %s", e.StatusCode, http.StatusText(e.StatusCode), e.Message)
}

// NewClient creates a client with the given authorization token.
func NewClient(id, secret string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		id:             id,
		secret:         secret,
		httpClient:     &http.Client{Timeout: time.Second * 60},
		tokenExpiresAt: time.Now(),

		// Set default values, which may be overridden by user options.
		baseURL:        apiURL,
		accessTokenURL: accessTokenURL,
		rateLimiter:    rate.NewLimiter(rate.Every(time.Second), 100),
	}

	// Apply options using the functional option pattern.
	for _, opt := range opts {
		opt(c)
	}

	// Get the access token.
	if _, err := c.getAccessToken(); err != nil {
		return nil, err
	}

	return c, nil
}

// WithSandbox sets the baseURL to the default sandbox URL.
func WithDefaultSandbox() ClientOption {
	return func(client *Client) {
		client.baseURL = sandboxURL
		client.accessTokenURL = sandboxTokenURL
	}
}

// WithBaseURL sets the baseURL for a new IEX Client.
func WithBaseURL(baseURL string) ClientOption {
	return func(client *Client) {
		client.baseURL = baseURL
	}
}

// WithRateLimiter sets the rate limiter.
func WithRateLimiter(duration time.Duration, numRequests int) ClientOption {
	return func(client *Client) {
		client.rateLimiter = rate.NewLimiter(rate.Every(duration), numRequests)
	}
}

// GetJSON gets the JSON data from the given endpoint.
func (c *Client) GetJSON(ctx context.Context, endpoint string, v any) error {
	u, err := c.url(endpoint, map[string]string{"token": c.accessToken})
	if err != nil {
		return err
	}
	return c.FetchURLToJSON(ctx, u, v)
}

// GetJSONWithQueryParams gets the JSON data from the given endpoint with the
// query parameters attached.
func (c *Client) GetJSONWithQueryParams(ctx context.Context,
	endpoint string, queryParams map[string]string, v interface{}) error {
	queryParams["token"] = c.accessToken
	u, err := c.url(endpoint, queryParams)
	if err != nil {
		return err
	}
	return c.FetchURLToJSON(ctx, u, v)
}

// Fetches JSON content from the given URL and unmarshals it into `v`.
func (c *Client) FetchURLToJSON(ctx context.Context, u *url.URL, v any) error {
	data, err := c.getBytes(ctx, u.String())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// GetJSONWithoutToken gets the JSON data from the given endpoint without
// adding a token to the URL.
func (c *Client) GetJSONWithoutToken(ctx context.Context, endpoint string, v any) error {
	u, err := c.url(endpoint, nil)
	if err != nil {
		return err
	}
	return c.FetchURLToJSON(ctx, u, v)
}

// GetBytes gets the data from the given endpoint.
func (c *Client) GetBytes(ctx context.Context, endpoint string) ([]byte, error) {
	u, err := c.url(endpoint, map[string]string{"token": c.accessToken})
	if err != nil {
		return nil, err
	}
	return c.getBytes(ctx, u.String())
}

// GetFloat64 gets the number from the given endpoint.
func (c *Client) GetFloat64(ctx context.Context, endpoint string) (float64, error) {
	b, err := c.GetBytes(ctx, endpoint)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(string(b), 64)
}

func (c *Client) getBytes(ctx context.Context, address string) ([]byte, error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return []byte{}, err
	}
	err = c.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// Even if GET didn't return an error, check the status code to make sure
	// everything was ok.
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		msg := ""

		if err == nil {
			msg = string(b)
		}

		return []byte{}, Error{StatusCode: resp.StatusCode, Message: msg}
	}
	return io.ReadAll(resp.Body)
}

// Returns a URL object that points to the endpoint with optional query parameters.
func (c *Client) url(endpoint string, queryParams map[string]string) (*url.URL, error) {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}
	return u, nil
}
