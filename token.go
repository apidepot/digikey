// Copyright (c) 2025 The digikey developers. All rights reserved.
// Project site: https://github.com/apidepot/digikey
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package digikey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// accessToken provides the response for a successful access token request.
// ExpiresIn is in seconds.
type accessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
	Type      string `json:"token_type"`
}

// getAccessToken returns the current access token or refreshes the access
// token using the client ID and client secret.
func (c *Client) getAccessToken() (string, error) {
	c.mu.RLock()
	if time.Now().Before(c.tokenExpiresAt) {
		token := c.accessToken
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	// Token is expred, so refresh.
	return c.refreshToken()
}

func (c *Client) refreshToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	requestBody := struct {
		ID     string `json:"client_id"`
		Secret string `json:"client_secret"`
		Type   string `json:"grant_type"`
	}{
		ID:     c.id,
		Secret: c.secret,
		Type:   grantType,
	}
	data, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling access token body: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.accessTokenURL,
		"application/x-www-form-urlencoded",
		bytes.NewReader(data),
	)
	if err != nil {
		return "", fmt.Errorf("error in post request for new access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"bad status code (%d) from post for new access token: %s",
			resp.StatusCode,
			string(errorBody),
		)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	accessToken := accessToken{}
	if err := json.Unmarshal(responseBody, &accessToken); err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	// Remove one second from the time to expriration to be safe.
	c.accessToken = accessToken.Token
	c.tokenType = accessToken.Type
	c.tokenExpiresAt = time.Now().Add(time.Duration(accessToken.ExpiresIn - 1))

	return c.accessToken, nil

}
