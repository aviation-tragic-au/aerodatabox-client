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
)

// ── EndpointHealthCheck ───────────────────────────────────────────────────────

func TestEndpointHealthCheck_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointHealthCheck
		want     string
	}{
		{EndpointHealthCheckFeedStatusGeneral, "/health/services/feeds/{service}"},
		{EndpointHealthCheckFeedStatusData, "/health/services/airports/{icao}/feeds"},
		{EndpointHealthCheckAirports, "/health/services/feeds/{service}/airports"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String(), "endpoint %d", tc.endpoint)
	}
}

func TestEndpointHealthCheck_Replace(t *testing.T) {
	t.Run("feed status general substitutes service token", func(t *testing.T) {
		got, err := EndpointHealthCheckFeedStatusGeneral.Replace(map[string]string{
			"service": "FlightSchedules",
		})
		require.NoError(t, err)
		assert.Equal(t, "/health/services/feeds/FlightSchedules", got)
	})

	t.Run("feed status data substitutes icao token", func(t *testing.T) {
		got, err := EndpointHealthCheckFeedStatusData.Replace(map[string]string{
			"icao": "EGLL",
		})
		require.NoError(t, err)
		assert.Equal(t, "/health/services/airports/EGLL/feeds", got)
	})

	t.Run("airports substitutes service token", func(t *testing.T) {
		got, err := EndpointHealthCheckAirports.Replace(map[string]string{
			"service": "FlightSchedules",
		})
		require.NoError(t, err)
		assert.Equal(t, "/health/services/feeds/FlightSchedules/airports", got)
	})

	t.Run("missing token returns error naming it", func(t *testing.T) {
		_, err := EndpointHealthCheckFeedStatusGeneral.Replace(map[string]string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "service")
	})
}

// ── GetDataFeedStatus ─────────────────────────────────────────────────────────

func TestGetDataFeedStatus(t *testing.T) {
	t.Run("calls general feed status endpoint without withHttpCode", func(t *testing.T) {
		expected := &model.FeedServiceStatusContract{
			Service: model.FeedServiceEnumFlightSchedules,
			Status:  model.FeedServiceStatusOK,
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/health/services/feeds/FlightSchedules", r.URL.Path)
			assert.Empty(t, r.URL.Query().Get("withHttpCode"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetDataFeedStatus("FlightSchedules", false)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, model.FeedServiceEnumFlightSchedules, got.Service)
		assert.Equal(t, model.FeedServiceStatusOK, got.Status)
	})

	t.Run("appends withHttpCode=true when requested", func(t *testing.T) {
		expected := &model.FeedServiceStatusContract{}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "true", r.URL.Query().Get("withHttpCode"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedStatus("FlightSchedules", true)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedStatus("FlightSchedules", false)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedStatus("FlightSchedules", false)
		require.Error(t, err)
	})
}

// ── GetDataFeedStatusByICAO ───────────────────────────────────────────────────

func TestGetDataFeedStatusByICAO(t *testing.T) {
	t.Run("returns feed status for airport", func(t *testing.T) {
		expected := &model.AirportFeedServiceStatusContract{
			FlightSchedulesFeed: model.FeedServiceStatusContract{
				Service: model.FeedServiceEnumFlightSchedules,
				Status:  model.FeedServiceStatusOK,
			},
			FlightLiveUpdates: model.FeedServiceStatusContract{
				Service: model.FeedServiceEnumFlightLiveUpdates,
				Status:  model.FeedServiceStatusDegraded,
			},
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/health/services/airports/EGLL/feeds", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetDataFeedStatusByICAO("EGLL")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, model.FeedServiceStatusOK, got.FlightSchedulesFeed.Status)
		assert.Equal(t, model.FeedServiceStatusDegraded, got.FlightLiveUpdates.Status)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedStatusByICAO("ZZZZ")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedStatusByICAO("EGLL")
		require.Error(t, err)
	})
}

// ── GetDataFeedSupportedAirports ──────────────────────────────────────────────

func TestGetDataFeedSupportedAirports(t *testing.T) {
	t.Run("returns airport list for service", func(t *testing.T) {
		expected := &model.StringCollectionContract{
			Count: 3,
			Items: []string{"EGLL", "KJFK", "YSSY"},
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/health/services/feeds/FlightSchedules/airports", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetDataFeedSupportedAirports("FlightSchedules")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, int64(3), got.Count)
		assert.Equal(t, []string{"EGLL", "KJFK", "YSSY"}, got.Items)
	})

	t.Run("returns empty collection", func(t *testing.T) {
		expected := &model.StringCollectionContract{Count: 0, Items: []string{}}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetDataFeedSupportedAirports("AdsbUpdates")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Empty(t, got.Items)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedSupportedAirports("FlightSchedules")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetDataFeedSupportedAirports("FlightSchedules")
		require.Error(t, err)
	})
}
