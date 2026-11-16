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
	"time"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── EndpointStatistical ───────────────────────────────────────────────────────

func TestEndpointStatistical_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointStatistical
		want     string
	}{
		{EndpointStatisticalAirportDelays, "/airports/{codeType}/{code}/delays"},
		{EndpointStatisticalAirportDelaysByDate, "/airports/{codeType}/{code}/delays/{dateLocal}"},
		{EndpointStatisticalAirportDelaysByDateRange, "/airports/{codeType}/{code}/delays/{dateFrom}/{dateTo}"},
		{EndpointStatisticalAirportRouteStatistics, "/airports/{codeType}/{code}/stats/routes/daily"},
		{EndpointStatisticalAirportRouteStatisticsByDate, "/airports/{codeType}/{code}/stats/routes/daily/{dateLocal}"},
		{EndpointStatisticalGlobalDelaysCurrent, "/airports/delays"},
		{EndpointStatisticalGlobalDelaysByDate, "/airports/delays/{dateUtc}"},
		{EndpointStatisticalFlightDelayStatistics, "/flights/{number}/delays"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String(), "endpoint %d", tc.endpoint)
	}
}

func TestEndpointStatistical_Replace(t *testing.T) {
	t.Run("airport delays all tokens substituted", func(t *testing.T) {
		got, err := EndpointStatisticalAirportDelays.Replace(map[string]string{
			"codeType": "icao",
			"code":     "EGLL",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/icao/EGLL/delays", got)
	})

	t.Run("airport delays by date all tokens substituted", func(t *testing.T) {
		got, err := EndpointStatisticalAirportDelaysByDate.Replace(map[string]string{
			"codeType":  "iata",
			"code":      "LHR",
			"dateLocal": "2026-01-15T00:00:00",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/iata/LHR/delays/2026-01-15T00:00:00", got)
	})

	t.Run("airport delays by date range all tokens substituted", func(t *testing.T) {
		got, err := EndpointStatisticalAirportDelaysByDateRange.Replace(map[string]string{
			"codeType": "icao",
			"code":     "EGLL",
			"dateFrom": "2026-01-01T00:00:00",
			"dateTo":   "2026-01-07T00:00:00",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/icao/EGLL/delays/2026-01-01T00:00:00/2026-01-07T00:00:00", got)
	})

	t.Run("flight delay statistics all tokens substituted", func(t *testing.T) {
		got, err := EndpointStatisticalFlightDelayStatistics.Replace(map[string]string{
			"number": "BA123",
		})
		require.NoError(t, err)
		assert.Equal(t, "/flights/BA123/delays", got)
	})

	t.Run("missing token returns error naming it", func(t *testing.T) {
		_, err := EndpointStatisticalAirportDelays.Replace(map[string]string{
			"codeType": "icao",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "code")
	})

	t.Run("no-token endpoint accepts empty map", func(t *testing.T) {
		got, err := EndpointStatisticalGlobalDelaysCurrent.Replace(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "/airports/delays", got)
	})

	t.Run("extra keys in map are ignored", func(t *testing.T) {
		got, err := EndpointStatisticalAirportDelays.Replace(map[string]string{
			"codeType": "iata",
			"code":     "JFK",
			"extra":    "ignored",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/iata/JFK/delays", got)
	})
}

// ── GetAirportDelays ──────────────────────────────────────────────────────────

func TestGetAirportDelays(t *testing.T) {
	t.Run("without date returns current delays", func(t *testing.T) {
		expected := &model.AirportDelayContract{
			AirportICAO: "EGLL",
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/delays", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportDelays("EGLL", "icao", nil)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "EGLL", got.AirportICAO)
	})

	t.Run("with date uses dateLocal endpoint", func(t *testing.T) {
		date := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		expected := &model.AirportDelayContract{
			AirportICAO: "EGLL",
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/delays/2026-01-15T00:00:00", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportDelays("EGLL", "icao", &date)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "EGLL", got.AirportICAO)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"airport not found"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportDelays("ZZZ", "icao", nil)
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

		_, err := newTestClient(srv.URL).GetAirportDelays("EGLL", "icao", nil)
		require.Error(t, err)
	})
}

// ── GetAirportDelaysInRange ───────────────────────────────────────────────────

func TestGetAirportDelaysInRange(t *testing.T) {
	t.Run("returns delays within date range", func(t *testing.T) {
		dateFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

		expected := []*model.AirportDelayContract{
			{AirportICAO: "EGLL"},
			{AirportICAO: "EGLL"},
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/delays/2026-01-01T00:00:00/2026-01-07T00:00:00", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportDelaysInRange("EGLL", "icao", dateFrom, dateTo)
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "EGLL", got[0].AirportICAO)
	})

	t.Run("returns empty list", func(t *testing.T) {
		dateFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("[]"))
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportDelaysInRange("EGLL", "icao", dateFrom, dateTo)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		dateFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportDelaysInRange("EGLL", "icao", dateFrom, dateTo)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		dateFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		dateTo := time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportDelaysInRange("EGLL", "icao", dateFrom, dateTo)
		require.Error(t, err)
	})
}

// ── GetAirportRoutes ──────────────────────────────────────────────────────────

func TestGetAirportRoutes(t *testing.T) {
	t.Run("returns route statistics", func(t *testing.T) {
		iata := "BA"
		expected := &model.DailyRouteStatContract{
			Routes: []model.DailyRouteStatRecordContract{
				{
					AverageDailyFlights: 3.5,
					Operations: []model.FlightAirlineContract{
						{Name: "British Airways", IATA: &iata},
					},
				},
			},
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/iata/LHR/stats/routes/daily", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportRoutes("LHR", "iata")
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got.Routes, 1)
		assert.Equal(t, 3.5, got.Routes[0].AverageDailyFlights)
		assert.Equal(t, "British Airways", got.Routes[0].Operations[0].Name)
	})

	t.Run("returns empty routes", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"routes":[]}`))
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportRoutes("LHR", "iata")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Empty(t, got.Routes)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"forbidden"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportRoutes("LHR", "iata")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportRoutes("LHR", "iata")
		require.Error(t, err)
	})
}

// ── GetAirportRoutesByDate ────────────────────────────────────────────────────

func TestGetAirportRoutesByDate(t *testing.T) {
	t.Run("returns route statistics for a specific date", func(t *testing.T) {
		date := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		expected := &model.DailyRouteStatContract{
			Routes: []model.DailyRouteStatRecordContract{
				{AverageDailyFlights: 2.0},
			},
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/stats/routes/daily/2026-03-10T00:00:00", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportRoutesByDate("EGLL", "icao", date)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got.Routes, 1)
		assert.Equal(t, 2.0, got.Routes[0].AverageDailyFlights)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		date := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"rate limit exceeded"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportRoutesByDate("EGLL", "icao", date)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		date := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportRoutesByDate("EGLL", "icao", date)
		require.Error(t, err)
	})
}

// ── GetGlobalDelays ───────────────────────────────────────────────────────────

func TestGetGlobalDelays(t *testing.T) {
	t.Run("returns global delays", func(t *testing.T) {
		expected := &model.AirportDelayContract{
			AirportICAO: "GLOBAL",
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/delays", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetGlobalDelays()
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "GLOBAL", got.AirportICAO)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"internal server error"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetGlobalDelays()
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
		assert.Equal(t, "internal server error", apiErr.Message)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetGlobalDelays()
		require.Error(t, err)
	})
}

// ── GetGlobalDelaysByDate ─────────────────────────────────────────────────────

func TestGetGlobalDelaysByDate(t *testing.T) {
	t.Run("returns global delays for a specific date", func(t *testing.T) {
		date := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
		expected := &model.AirportDelayContract{
			AirportICAO: "GLOBAL",
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/delays/2026-06-01T12:00:00", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetGlobalDelaysByDate(date)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "GLOBAL", got.AirportICAO)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		date := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetGlobalDelaysByDate(date)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		date := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetGlobalDelaysByDate(date)
		require.Error(t, err)
	})
}

// ── GetFlightDelayStatistics ──────────────────────────────────────────────────

func TestGetFlightDelayStatistics(t *testing.T) {
	t.Run("returns flight delay statistics", func(t *testing.T) {
		expected := &model.FlightDelayContract{
			AirportICAO: "EGLL",
			MedianDelay: "PT5M",
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/BA123/delays", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightDelayStatistics("BA123")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "EGLL", got.AirportICAO)
		assert.Equal(t, "PT5M", got.MedianDelay)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"flight not found"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightDelayStatistics("XX999")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
		assert.Equal(t, "flight not found", apiErr.Message)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightDelayStatistics("BA123")
		require.Error(t, err)
	})
}
