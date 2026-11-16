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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"go.uber.org/zap"
)

const (
	// DefaultHost is AeroDataBox's RapidAPI host, used when Config.Host is
	// empty.
	DefaultHost = "aerodatabox.p.rapidapi.com"

	// DefaultTimeout bounds a request when Config.HTTPClient is unset.
	DefaultTimeout = 10 * time.Second

	// HeaderRapidAPIKey carries the RapidAPI key.
	HeaderRapidAPIKey = "X-RapidAPI-Key"
	// HeaderRapidAPIHost names the RapidAPI host the request is for.
	HeaderRapidAPIHost = "X-RapidAPI-Host"
)

// ErrNoData means AeroDataBox has nothing for the request: it answered 204
// or an empty body (it doesn't cover every airport or flight), or 404. A
// 404 is also an *model.APIError, for its message.
var ErrNoData = model.ErrNoData

// Config configures a Client. Only a key (Key or KeyFunc) is required. The
// mapstructure tags let a host decode it straight from viper or similar.
type Config struct {
	// Host is the RapidAPI host, sent as X-RapidAPI-Host. Empty means
	// DefaultHost.
	Host string `json:"host" mapstructure:"host"`

	// URL is the API base. Empty means https://<Host>. Set it to reach
	// AeroDataBox some other way (its own API.Market endpoint, a proxy,
	// a test server).
	URL string `json:"url" mapstructure:"url"`

	// Key is the RapidAPI key.
	Key string `json:"key" mapstructure:"key"`

	// KeyFunc, when set, is called for every request in place of Key, so
	// a host can change the key (from a settings page, say) without
	// building a new client.
	KeyFunc func() string `json:"-" mapstructure:"-"`

	// HTTPClient sends the requests. Nil means a client with
	// DefaultTimeout.
	HTTPClient *http.Client `json:"-" mapstructure:"-"`
}

// Client is a client for the AeroDataBox API. It is safe for concurrent
// use, as long as URL and Host aren't changed while it's in use.
type Client struct {
	logger     *zap.Logger
	httpClient *http.Client
	key        func() string

	// URL is the API base requests are sent to.
	URL string
	// Host is sent as X-RapidAPI-Host.
	Host string
}

// NewClient creates an AeroDataBox client. A nil logger logs nothing.
func NewClient(logger *zap.Logger, cfg Config) (*Client, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	host := cfg.Host
	if host == "" {
		host = DefaultHost
	}

	base := strings.TrimRight(cfg.URL, "/")
	if base == "" {
		base = "https://" + host
	}
	if _, err := url.Parse(base); err != nil {
		return nil, fmt.Errorf("aerodatabox: invalid URL: %w", err)
	}

	key := cfg.KeyFunc
	if key == nil {
		static := cfg.Key
		key = func() string { return static }
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultTimeout}
	}

	return &Client{
		logger:     logger.Named("aerodatabox"),
		httpClient: httpClient,
		key:        key,
		URL:        base,
		Host:       host,
	}, nil
}

func (c *Client) getRequest(ctx context.Context, endpoint string) (*http.Request, error) {
	return c.newRequest(ctx, http.MethodGet, endpoint, nil)
}

func (c *Client) postRequest(ctx context.Context, endpoint string, payload any) (*http.Request, error) {
	return c.newRequest(ctx, http.MethodPost, endpoint, payload)
}

func (c *Client) deleteRequest(ctx context.Context, endpoint string) (*http.Request, error) {
	return c.newRequest(ctx, http.MethodDelete, endpoint, nil)
}

// newRequest builds a request with the RapidAPI headers, and a JSON body
// when payload isn't nil.
func (c *Client) newRequest(ctx context.Context, method, endpoint string, payload any) (*http.Request, error) {
	key := c.key()
	if key == "" {
		return nil, errors.New("aerodatabox: no API key set")
	}

	var body io.Reader = http.NoBody
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("aerodatabox: unable to encode request body: %w", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.URL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("aerodatabox: unable to build request: %w", err)
	}

	req.Header.Set(HeaderRapidAPIKey, key)
	req.Header.Set(HeaderRapidAPIHost, c.Host)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// executeRequest sends req and decodes a 2xx JSON body into result. With a
// nil result (a delete, say) any 2xx is success; otherwise a 204 or an
// empty body is ErrNoData. A non-2xx is an *model.APIError.
func (c *Client) executeRequest(req *http.Request, result any) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("aerodatabox: request failed: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("aerodatabox: unable to read response: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		apiErr := &model.APIError{StatusCode: res.StatusCode}
		if len(body) > 0 {
			if err := json.Unmarshal(body, apiErr); err != nil || apiErr.Message == "" {
				apiErr.Message = quote(body)
			}
		}
		c.logger.Debug(
			"API error response",
			zap.String("path", req.URL.Path),
			zap.Int("status", res.StatusCode),
			zap.String("message", apiErr.Message),
		)
		return apiErr
	}

	if result == nil {
		return nil
	}
	if res.StatusCode == http.StatusNoContent || len(bytes.TrimSpace(body)) == 0 {
		return ErrNoData
	}

	if err := json.Unmarshal(body, result); err != nil {
		c.logger.Debug("unable to parse response", zap.String("path", req.URL.Path), zap.ByteString("body", body))
		return fmt.Errorf("aerodatabox: unable to parse response: %w", err)
	}

	return nil
}

// quote returns an error response body fit for an error message: trimmed,
// and cut short if it's long.
func quote(body []byte) string {
	const maxErrorBody = 512
	s := strings.TrimSpace(string(body))
	if len(s) > maxErrorBody {
		s = s[:maxErrorBody] + "..."
	}
	return s
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
// with corresponding values from the provided map, path-escaped. An error
// is returned if any token present in the template has no matching key in
// the map.
func replaceTokens(template string, tokens map[string]string) (string, error) {
	var missing []string

	result := tokenPattern.ReplaceAllStringFunc(template, func(match string) string {
		key := match[1 : len(match)-1] // strip surrounding { }
		if val, ok := tokens[key]; ok {
			return url.PathEscape(val)
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
