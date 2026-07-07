package client

//    Copyright 2026 Dean Stalker
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/deanstalker-stream/aerodatabox-client/model"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	Namespace          = "feed.aerodatabox"
	HeaderRapidAPIKey  = "X-RapidAPI-Key"
	HeaderRapidAPIHost = "X-RapidAPI-Host"
)

// BindEnvs registers environment variable mappings for this feed's config namespace.
func BindEnvs(v *viper.Viper) {
	_ = v.BindEnv("feed.aerodatabox.host", "FEED_AERODATABOX_HOST")
	_ = v.BindEnv("feed.aerodatabox.key", "FEED_AERODATABOX_KEY")
}

// Client represents a client for the AeroDataBox API.
type Client struct {
	logger *zap.Logger

	URL  string
	Host string
	Key  string
}

// NewClient creates a new AeroDataBox client.
func NewClient(logger *zap.Logger, cfg *model.Config) (*Client, error) {
	return &Client{
		logger: logger.Named(Namespace),
		URL:    fmt.Sprintf("https://%s", cfg.Host),
		Host:   cfg.Host,
		Key:    cfg.Key,
	}, nil
}

func (c *Client) getRequest(endpoint string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.URL, endpoint)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HeaderRapidAPIKey, c.Key)
	req.Header.Add(HeaderRapidAPIHost, c.Host)

	return req, nil
}

func (c *Client) postRequest(endpoint string, payload interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.URL, endpoint)

	var body io.Reader = http.NoBody
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HeaderRapidAPIKey, c.Key)
	req.Header.Add(HeaderRapidAPIHost, c.Host)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) deleteRequest(endpoint string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.URL, endpoint)
	req, err := http.NewRequest("DELETE", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add(HeaderRapidAPIKey, c.Key)
	req.Header.Add(HeaderRapidAPIHost, c.Host)

	return req, nil
}

func (c *Client) executeRequest(req *http.Request, result interface{}) error {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		apiErr := &model.APIError{StatusCode: res.StatusCode}
		if len(body) > 0 {
			if err := json.Unmarshal(body, apiErr); err != nil {
				apiErr.Message = string(body)
			}
		}
		c.logger.Error(
			"API error response",
			zap.Int("status", res.StatusCode),
			zap.String("message", apiErr.Message),
		)
		return apiErr
	}

	if err := json.Unmarshal(body, result); err != nil {
		c.logger.Error("unable to parse response", zap.Error(err), zap.ByteString("body", body))
		return err
	}

	return nil
}

func addQueryParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
}

var tokenPattern = regexp.MustCompile(`\{([^}]+)\}`)

// replaceTokens substitutes all {token} placeholders in the given template
// with corresponding values from the provided map. An error is returned if
// any token present in the template has no matching key in the map.
func replaceTokens(template string, tokens map[string]string) (string, error) {
	var missing []string

	result := tokenPattern.ReplaceAllStringFunc(template, func(match string) string {
		key := match[1 : len(match)-1] // strip surrounding { }
		if val, ok := tokens[key]; ok {
			return val
		}
		missing = append(missing, key)
		return match
	})

	if len(missing) > 0 {
		return "", fmt.Errorf("endpoint %q: missing token values for: %s", template, strings.Join(missing, ", "))
	}

	return result, nil
}

// appendQuery appends pre-encoded query parameters to a base URL string.
// Returns base unchanged if encoded is empty.
func appendQuery(base, encoded string) string {
	if encoded != "" {
		return base + "?" + encoded
	}
	return base
}

func (c *Client) validateCode(codeType, code string) error {
	codeType = strings.ToLower(codeType)
	switch codeType {
	case "iata":
		if len(code) != 3 {
			return fmt.Errorf("invalid IATA code: %s", code)
		}
	case "icao":
		if len(code) != 4 {
			return fmt.Errorf("invalid ICAO code: %s", code)
		}
	default:
		return fmt.Errorf("invalid code type: %s", codeType)
	}
	return nil
}
