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
	"testing"
	"time"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── EndpointMisc ─────────────────────────────────────────────────────────────

func TestEndpointMisc_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointMisc
		want     string
	}{
		{EndpointMiscAirportCurrentTime, "/airports/{codeType}/{code}/time/local"},
		{EndpointMiscAirportSolarTime, "/airports/{codeType}/{code}/time/solar"},
		{EndpointMiscAirportSolarTimeByDate, "/airports/{codeType}/{code}/time/solar/{dateLocal}"},
		{EndpointMiscFlightTime, "/airports/{codeType}/{codeFrom}/distance-time/{codeTo}"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String(), "endpoint %d", tc.endpoint)
	}
}

func TestEndpointMisc_Replace(t *testing.T) {
	t.Run("airport current time all tokens substituted", func(t *testing.T) {
		got, err := EndpointMiscAirportCurrentTime.Replace(map[string]string{
			"codeType": "iata",
			"code":     "JFK",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/iata/JFK/time/local", got)
	})

	t.Run("airport solar time by date all tokens substituted", func(t *testing.T) {
		got, err := EndpointMiscAirportSolarTimeByDate.Replace(map[string]string{
			"codeType":  "iata",
			"code":      "JFK",
			"dateLocal": "2026-06-21T12:00",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/iata/JFK/time/solar/2026-06-21T12:00", got)
	})

	t.Run("flight time all tokens substituted", func(t *testing.T) {
		got, err := EndpointMiscFlightTime.Replace(map[string]string{
			"codeType": "iata",
			"codeFrom": "JFK",
			"codeTo":   "LHR",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/iata/JFK/distance-time/LHR", got)
	})

	t.Run("missing token returns error naming it", func(t *testing.T) {
		_, err := EndpointMiscAirportCurrentTime.Replace(map[string]string{
			"codeType": "iata",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "code")
	})

	t.Run("extra keys in map are ignored", func(t *testing.T) {
		got, err := EndpointMiscAirportSolarTime.Replace(map[string]string{
			"codeType": "icao",
			"code":     "EGLL",
			"extra":    "ignored",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/icao/EGLL/time/solar", got)
	})
}

// ── GetAirportCurrentTime ─────────────────────────────────────────────────────

func TestGetAirportCurrentTime(t *testing.T) {
	t.Run("returns current time for valid IATA code", func(t *testing.T) {
		expected := &model.AirportLocalTimeContract{TimeZoneID: "America/New_York"}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/iata/JFK/time/local", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportCurrentTime(context.Background(), "JFK", "iata")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "America/New_York", got.TimeZoneID)
	})

	t.Run("returns current time for valid ICAO code", func(t *testing.T) {
		expected := &model.AirportLocalTimeContract{TimeZoneID: "Europe/London"}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/time/local", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportCurrentTime(context.Background(), "EGLL", "icao")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "Europe/London", got.TimeZoneID)
	})

	t.Run("invalid IATA code returns error without making request", func(t *testing.T) {
		got, err := newTestClient("").GetAirportCurrentTime(context.Background(), "TOOLONG", "iata")
		require.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "IATA")
	})

	t.Run("invalid ICAO code returns error without making request", func(t *testing.T) {
		got, err := newTestClient("").GetAirportCurrentTime(context.Background(), "EG", "icao")
		require.Error(t, err)
		assert.Nil(t, got)
		assert.Contains(t, err.Error(), "ICAO")
	})

	t.Run("invalid code type returns error without making request", func(t *testing.T) {
		got, err := newTestClient("").GetAirportCurrentTime(context.Background(), "JFK", "local")
		require.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportCurrentTime(context.Background(), "JFK", "iata")
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

		_, err := newTestClient(srv.URL).GetAirportCurrentTime(context.Background(), "JFK", "iata")
		require.Error(t, err)
	})
}

// ── GetAirportSolarTime ───────────────────────────────────────────────────────

func TestGetAirportSolarTime(t *testing.T) {
	t.Run("returns solar state for valid ICAO code", func(t *testing.T) {
		expected := &model.SolarStateContract{}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/time/solar", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportSolarTime(context.Background(), "EGLL", "icao")
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("invalid IATA code returns error without making request", func(t *testing.T) {
		_, err := newTestClient("").GetAirportSolarTime(context.Background(), "JK", "iata")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "IATA")
	})

	t.Run("invalid ICAO code returns error without making request", func(t *testing.T) {
		_, err := newTestClient("").GetAirportSolarTime(context.Background(), "TOOLONG", "icao")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ICAO")
	})

	t.Run("invalid code type returns error without making request", func(t *testing.T) {
		_, err := newTestClient("").GetAirportSolarTime(context.Background(), "JFK", "unknown")
		require.Error(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportSolarTime(context.Background(), "JFK", "iata")
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

		_, err := newTestClient(srv.URL).GetAirportSolarTime(context.Background(), "JFK", "iata")
		require.Error(t, err)
	})
}

// ── GetAirportSolarTimeByDate ─────────────────────────────────────────────────

func TestGetAirportSolarTimeByDate(t *testing.T) {
	date := time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC)

	t.Run("includes formatted date in path", func(t *testing.T) {
		expected := &model.SolarStateContract{}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/time/solar/2026-06-21T00:00", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportSolarTimeByDate(context.Background(), "EGLL", "icao", date)
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("uses 2006-01-02T15:04 date format", func(t *testing.T) {
		noon := time.Date(2026, 12, 25, 14, 30, 0, 0, time.UTC)
		expected := &model.SolarStateContract{}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/iata/JFK/time/solar/2026-12-25T14:30", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportSolarTimeByDate(context.Background(), "JFK", "iata", noon)
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("invalid IATA code returns error without making request", func(t *testing.T) {
		_, err := newTestClient("").GetAirportSolarTimeByDate(context.Background(), "JK", "iata", date)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "IATA")
	})

	t.Run("invalid ICAO code returns error without making request", func(t *testing.T) {
		_, err := newTestClient("").GetAirportSolarTimeByDate(context.Background(), "EG", "icao", date)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ICAO")
	})

	t.Run("invalid code type returns error without making request", func(t *testing.T) {
		_, err := newTestClient("").GetAirportSolarTimeByDate(context.Background(), "JFK", "local", date)
		require.Error(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportSolarTimeByDate(context.Background(), "JFK", "iata", date)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportSolarTimeByDate(context.Background(), "JFK", "iata", date)
		require.Error(t, err)
	})
}

// ── GetGreatCircleDistance ────────────────────────────────────────────────────

func TestGetGreatCircleDistance(t *testing.T) {
	t.Run("returns distance for IATA codes", func(t *testing.T) {
		expected := &model.AirportDistanceTimeContract{ApproximateFlightTime: "PT7H00M"}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/iata/JFK/distance-time/LHR", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetGreatCircleDistance(context.Background(), "iata", "JFK", "LHR", nil, nil)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "PT7H00M", got.ApproximateFlightTime)
	})

	t.Run("returns distance for ICAO codes", func(t *testing.T) {
		expected := &model.AirportDistanceTimeContract{ApproximateFlightTime: "PT7H15M"}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/KJFK/distance-time/EGLL", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetGreatCircleDistance(context.Background(), "icao", "KJFK", "EGLL", nil, nil)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "PT7H15M", got.ApproximateFlightTime)
	})

	t.Run("sends flightTimeModel and aircraftName as query params", func(t *testing.T) {
		flightModel := "ML01"
		aircraft := "Boeing 737"
		expected := &model.AirportDistanceTimeContract{}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "ML01", r.URL.Query().Get("flightTimeModel"))
			assert.Equal(t, "Boeing 737", r.URL.Query().Get("aircraftName"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetGreatCircleDistance(context.Background(), "iata", "JFK", "LHR", &aircraft, &flightModel)
		require.NoError(t, err)
	})

	t.Run("codeType is case-insensitive", func(t *testing.T) {
		expected := &model.AirportDistanceTimeContract{}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/iata/JFK/distance-time/LHR", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetGreatCircleDistance(context.Background(), "IATA", "JFK", "LHR", nil, nil)
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("invalid IATA codeFrom returns error", func(t *testing.T) {
		_, err := newTestClient("").GetGreatCircleDistance(context.Background(), "iata", "TOOLONG", "LHR", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "IATA")
	})

	t.Run("invalid IATA codeTo returns error", func(t *testing.T) {
		_, err := newTestClient("").GetGreatCircleDistance(context.Background(), "iata", "JFK", "TOOLONG", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "IATA")
	})

	t.Run("invalid ICAO codeFrom returns error", func(t *testing.T) {
		_, err := newTestClient("").GetGreatCircleDistance(context.Background(), "icao", "EG", "EGLL", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ICAO")
	})

	t.Run("invalid ICAO codeTo returns error", func(t *testing.T) {
		_, err := newTestClient("").GetGreatCircleDistance(context.Background(), "icao", "KJFK", "EG", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ICAO")
	})

	t.Run("invalid code type returns error", func(t *testing.T) {
		_, err := newTestClient("").GetGreatCircleDistance(context.Background(), "local", "JFK", "LHR", nil, nil)
		require.Error(t, err)
	})

	t.Run("invalid flightTimeModel returns error", func(t *testing.T) {
		bad := "InvalidModel"
		_, err := newTestClient("").GetGreatCircleDistance(context.Background(), "iata", "JFK", "LHR", nil, &bad)
		require.Error(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetGreatCircleDistance(context.Background(), "iata", "JFK", "LHR", nil, nil)
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

		_, err := newTestClient(srv.URL).GetGreatCircleDistance(context.Background(), "iata", "JFK", "LHR", nil, nil)
		require.Error(t, err)
	})
}
