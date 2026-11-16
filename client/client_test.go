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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	cfg := &model.Config{Host: "api.aerodatabox.com", Key: "test-key"}

	client, err := NewClient(zap.NewNop(), cfg)

	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.aerodatabox.com", client.URL)
	assert.Equal(t, "api.aerodatabox.com", client.Host)
	assert.Equal(t, "test-key", client.Key)
	assert.NotNil(t, client.logger)
}

func TestAddQueryParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", http.NoBody)

	addQueryParams(req, map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})

	query := req.URL.Query()
	assert.Equal(t, "value1", query.Get("key1"))
	assert.Equal(t, "value2", query.Get("key2"))
	assert.Equal(t, "value3", query.Get("key3"))
}

func TestGetRequest(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.aerodatabox.com",
		Host:   "api.aerodatabox.com",
		Key:    "test-key-123",
	}

	req, err := client.getRequest("/airports/iata/JFK")

	require.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://api.aerodatabox.com/airports/iata/JFK", req.URL.String())
	assert.Equal(t, "test-key-123", req.Header.Get(HeaderRapidAPIKey))
	assert.Equal(t, "api.aerodatabox.com", req.Header.Get(HeaderRapidAPIHost))
}

func TestGetRequest_InvalidURL(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "ht!tp://[invalid",
		Host:   "api.aerodatabox.com",
		Key:    "test-key",
	}

	req, err := client.getRequest("/airports/iata/JFK")

	assert.Error(t, err)
	assert.Nil(t, req)
}

func TestExecuteRequest_Unmarshal(t *testing.T) {
	client := &Client{logger: zap.NewNop()}

	t.Run("successful unmarshal to timezone", func(t *testing.T) {
		expected := &model.Timezone{TimeZone: "America/New_York"}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "America/New_York", result.TimeZone)
	})

	t.Run("unmarshal to flights array", func(t *testing.T) {
		expected := []*model.Flight{{Number: "BA123", Status: "Scheduled"}}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result []*model.Flight
		err := client.executeRequest(req, &result)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "BA123", result[0].Number)
	})
}

func TestExecuteRequest_ErrorHandling(t *testing.T) {
	client := &Client{logger: zap.NewNop()}

	t.Run("400 with full ErrorContract body", func(t *testing.T) {
		details := "invalid param"
		body, _ := json.Marshal(&model.APIError{Message: "bad request", Details: &details})

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(body)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		require.Error(t, err)
		assert.Nil(t, result)

		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "bad request", apiErr.Message)
		require.NotNil(t, apiErr.Details)
		assert.Equal(t, "invalid param", *apiErr.Details)
	})

	t.Run("400 with message-only ErrorContract body", func(t *testing.T) {
		body, _ := json.Marshal(&model.APIError{Message: "not found"})

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(body)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		require.Error(t, err)

		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "not found", apiErr.Message)
		assert.Nil(t, apiErr.Details)
	})

	t.Run("503 with no body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		require.Error(t, err)

		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	})

	t.Run("204 No Content returns no error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestHeaderConstants(t *testing.T) {
	assert.Equal(t, "X-RapidAPI-Key", HeaderRapidAPIKey)
	assert.Equal(t, "X-RapidAPI-Host", HeaderRapidAPIHost)
}

func TestGetWeather_RequestStructure(t *testing.T) {
	client := &Client{
		logger: zap.NewNop(),
		URL:    "https://api.aerodatabox.com",
		Host:   "api.aerodatabox.com",
		Key:    "test-key",
	}

	req, err := client.getRequest("/airports/iata/JFK/weather")

	require.NoError(t, err)
	assert.Contains(t, req.URL.String(), "/airports/iata/JFK/weather")
}

func TestGetSingleAircraft(t *testing.T) {
	boolTrue := true

	t.Run("successful response with params", func(t *testing.T) {
		expected := &model.Airframe{
			Registration: "G-ABCD",
			HexICAO:      "400000",
			TypeName:     "Boeing 737",
		}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/aircrafts/reg/G-ABCD", r.URL.Path) //nolint:misspell
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := &Client{
			logger: zap.NewNop(),
			URL:    server.URL,
			Host:   "api.aerodatabox.com",
			Key:    "test-key",
		}

		result, err := c.GetSingleAircraft("reg", "G-ABCD", &boolTrue, &boolTrue)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "G-ABCD", result.Registration)
	})

	t.Run("successful response with nil params", func(t *testing.T) {
		expected := &model.Airframe{
			Registration: "N12345",
			HexICAO:      "A00001",
		}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/aircrafts/icao/A00001", r.URL.Path) //nolint:misspell
			assert.Empty(t, r.URL.Query().Get("withRegistrations"))
			assert.Empty(t, r.URL.Query().Get("withImage"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := &Client{
			logger: zap.NewNop(),
			URL:    server.URL,
			Host:   "api.aerodatabox.com",
			Key:    "test-key",
		}

		result, err := c.GetSingleAircraft("icao", "A00001", nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "N12345", result.Registration)
	})

	t.Run("invalid JSON response returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not valid json"))
		}))
		defer server.Close()

		c := &Client{
			logger: zap.NewNop(),
			URL:    server.URL,
			Host:   "api.aerodatabox.com",
			Key:    "test-key",
		}

		result, err := c.GetSingleAircraft("reg", "G-ABCD", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
