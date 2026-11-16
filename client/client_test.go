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
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		c, err := NewClient(nil, Config{Key: "k"})
		require.NoError(t, err)
		assert.Equal(t, "https://"+DefaultHost, c.URL)
		assert.Equal(t, DefaultHost, c.Host)
		assert.Equal(t, DefaultTimeout, c.httpClient.Timeout)
		assert.Equal(t, "k", c.key())
	})

	t.Run("host gives the URL", func(t *testing.T) {
		c, err := NewClient(nil, Config{Host: "api.example.test", Key: "k"})
		require.NoError(t, err)
		assert.Equal(t, "https://api.example.test", c.URL)
		assert.Equal(t, "api.example.test", c.Host)
	})

	t.Run("URL overrides the host's", func(t *testing.T) {
		c, err := NewClient(nil, Config{URL: "http://localhost:8080/", Key: "k"})
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080", c.URL)
		assert.Equal(t, DefaultHost, c.Host)
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := NewClient(nil, Config{URL: "http://bad\x7f", Key: "k"})
		require.Error(t, err)
	})

	t.Run("own HTTP client", func(t *testing.T) {
		hc := &http.Client{}
		c, err := NewClient(nil, Config{Key: "k", HTTPClient: hc})
		require.NoError(t, err)
		assert.Same(t, hc, c.httpClient)
	})
}

func TestKeyFuncReadPerRequest(t *testing.T) {
	var n atomic.Int32
	keys := []string{"first", "second"}
	c, err := NewClient(nil, Config{Key: "ignored", KeyFunc: func() string { return keys[n.Add(1)-1] }})
	require.NoError(t, err)

	for _, want := range keys {
		req, err := c.getRequest(context.Background(), "/x")
		require.NoError(t, err)
		assert.Equal(t, want, req.Header.Get(HeaderRapidAPIKey))
	}
}

func TestNoKey(t *testing.T) {
	c, err := NewClient(nil, Config{})
	require.NoError(t, err)

	_, err = c.GetAirportRunways(context.Background(), "YBBN", "icao")
	require.ErrorContains(t, err, "no API key set")
}

func TestCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := newTestClient(srv.URL).GetAirportRunways(ctx, "YBBN", "icao")
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestReplaceTokensEscapes(t *testing.T) {
	got, err := replaceTokens("/flights/{by}/{q}", map[string]string{"by": "callsign", "q": "A/B C"})
	require.NoError(t, err)
	assert.Equal(t, "/flights/callsign/A%2FB%20C", got)
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
	client := newTestClient("https://api.aerodatabox.com")

	req, err := client.getRequest(context.Background(), "/airports/iata/JFK")

	require.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://api.aerodatabox.com/airports/iata/JFK", req.URL.String())
	assert.Equal(t, "test-key", req.Header.Get(HeaderRapidAPIKey))
	assert.Equal(t, "api.aerodatabox.com", req.Header.Get(HeaderRapidAPIHost))
}

func TestGetRequest_InvalidURL(t *testing.T) {
	client := newTestClient("")
	client.URL = "http://bad\x7f"

	req, err := client.getRequest(context.Background(), "/airports/iata/JFK")

	assert.Error(t, err)
	assert.Nil(t, req)
}

func TestExecuteRequest_Unmarshal(t *testing.T) {
	client := newTestClient("")

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
	client := newTestClient("")

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

	t.Run("204 No Content is ErrNoData", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		require.ErrorIs(t, err, ErrNoData)
		assert.Nil(t, result)
	})

	t.Run("empty 200 is ErrNoData", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		require.ErrorIs(t, client.executeRequest(req, &result), ErrNoData)
	})

	t.Run("404 is an APIError and ErrNoData", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"Flight not found"}`))
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)

		require.ErrorIs(t, err, ErrNoData)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "Flight not found", apiErr.Message)
	})

	t.Run("400 is not ErrNoData", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		err := client.executeRequest(req, &result)
		require.Error(t, err)
		assert.NotErrorIs(t, err, ErrNoData)
	})

	t.Run("non-JSON error body is quoted", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("<html>bad gateway</html>"))
		}))
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, http.NoBody)
		var result *model.Timezone
		var apiErr *model.APIError
		require.ErrorAs(t, client.executeRequest(req, &result), &apiErr)
		assert.Equal(t, "<html>bad gateway</html>", apiErr.Message)
	})

	t.Run("204 with no result wanted is success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL, http.NoBody)
		assert.NoError(t, client.executeRequest(req, nil))
	})
}

func TestHeaderConstants(t *testing.T) {
	assert.Equal(t, "X-RapidAPI-Key", HeaderRapidAPIKey)
	assert.Equal(t, "X-RapidAPI-Host", HeaderRapidAPIHost)
}

func TestGetWeather_RequestStructure(t *testing.T) {
	client := newTestClient("https://api.aerodatabox.com")

	req, err := client.getRequest(context.Background(), "/airports/iata/JFK/weather")

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

		c := newTestClient(server.URL)

		result, err := c.GetSingleAircraft(context.Background(), "reg", "G-ABCD", &boolTrue, &boolTrue)

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

		c := newTestClient(server.URL)

		result, err := c.GetSingleAircraft(context.Background(), "icao", "A00001", nil, nil)

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

		c := newTestClient(server.URL)

		result, err := c.GetSingleAircraft(context.Background(), "reg", "G-ABCD", nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
