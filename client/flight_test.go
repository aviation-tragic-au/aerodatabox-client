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

	"github.com/deanstalker-stream/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── EndpointFlight ────────────────────────────────────────────────────────────

func TestEndpointFlight_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointFlight
		want     string
	}{
		{EndpointFlightStatus, "/flights/{searchBy}/{searchParameter}"},
		{EndpointFlightStatusByDate, "/flights/{searchBy}/{searchParameter}/{dateLocal}"},
		{EndpointFlightScheduleHistory, "/flights/{searchBy}/{searchParameter}/{dateFromLocal}/{dateToLocal}"},
		{EndpointFlightDepartureDates, "/flights/{searchBy}/{searchParameter}/dates"},
		{EndpointFlightDepartureDateRange, "/flights/{searchBy}/{searchParameter}/dates/{fromLocal}/{toLocal}"},
		{EndpointFlightFIDSLocalTime, "/flights/airports/{codeType}/{code}/{fromLocal}/{toLocal}"},
		{EndpointFlightFIDSRelativeTime, "/flights/airports/{codeType}/{code}"},
		{EndpointFlightSearchByTerm, "/flights/search/term"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String())
	}
}

func TestEndpointFlight_Replace(t *testing.T) {
	t.Run("status endpoint substitutes all tokens", func(t *testing.T) {
		got, err := EndpointFlightStatus.Replace(map[string]string{
			"searchBy":        "Number",
			"searchParameter": "BA123",
		})
		require.NoError(t, err)
		assert.Equal(t, "/flights/Number/BA123", got)
	})

	t.Run("missing token returns error", func(t *testing.T) {
		_, err := EndpointFlightStatus.Replace(map[string]string{"searchBy": "Number"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "searchParameter")
	})

	t.Run("no-token endpoint accepts empty map", func(t *testing.T) {
		got, err := EndpointFlightSearchByTerm.Replace(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "/flights/search/term", got)
	})

	t.Run("extra keys are ignored", func(t *testing.T) {
		got, err := EndpointFlightStatus.Replace(map[string]string{
			"searchBy":        "Reg",
			"searchParameter": "G-ABCD",
			"extra":           "ignored",
		})
		require.NoError(t, err)
		assert.Equal(t, "/flights/Reg/G-ABCD", got)
	})
}

// ── GetFlightStatus ───────────────────────────────────────────────────────────

func TestGetFlightStatus(t *testing.T) {
	t.Run("returns flight contracts with correct path and query params", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/Number/BA123", r.URL.Path)
			assert.Equal(t, "true", r.URL.Query().Get("withAircraftImage"))
			assert.Equal(t, "true", r.URL.Query().Get("withLocation"))
			assert.Equal(t, "false", r.URL.Query().Get("withFlightPlan"))
			assert.Equal(t, "Arrival", r.URL.Query().Get("dateLocalRole"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]*model.FlightContract{{Number: "BA123"}})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightStatus("Number", "BA123", model.FlightDirectionArrival, true, true, false)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "BA123", got[0].Number)
	})

	t.Run("rejects searchParameter longer than 10 chars", func(t *testing.T) {
		_, err := newTestClient("http://unused").GetFlightStatus("Number", "TOOLONGPARAM", model.FlightDirectionArrival, false, false, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds maximum length")
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"flight not found"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightStatus("Number", "XX999", model.FlightDirectionArrival, false, false, false)
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

		_, err := newTestClient(srv.URL).GetFlightStatus("Number", "BA123", model.FlightDirectionArrival, false, false, false)
		require.Error(t, err)
	})
}

// ── GetFlightStatusByDate ─────────────────────────────────────────────────────

func TestGetFlightStatusByDate(t *testing.T) {
	date := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

	t.Run("encodes date in path and sends query params", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/Number/BA123/2024-01-15T14:30", r.URL.Path)
			assert.Equal(t, "true", r.URL.Query().Get("withAircraftImage"))
			assert.Equal(t, "false", r.URL.Query().Get("withLocation"))
			assert.Equal(t, "false", r.URL.Query().Get("withFlightPlan"))
			assert.Equal(t, "Departure", r.URL.Query().Get("dateLocalRole"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]*model.FlightContract{{Number: "BA123"}})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightStatusByDate("Number", "BA123", date, model.FlightDirectionDeparture, true, false, false)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "BA123", got[0].Number)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightStatusByDate("Number", "BA123", date, model.FlightDirectionArrival, false, false, false)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})
}

// ── GetFlightScheduleHistory ──────────────────────────────────────────────────

func TestGetFlightScheduleHistory(t *testing.T) {
	dateFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2024, 1, 7, 23, 59, 0, 0, time.UTC)

	t.Run("returns history without dateLocalRole filter", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/Number/BA123/2024-01-01T00:00/2024-01-07T23:59", r.URL.Path)
			assert.Empty(t, r.URL.Query().Get("dateLocalRole"))
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]*model.FlightContract{{Number: "BA123"}})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightScheduleHistory("Number", "BA123", dateFrom, dateTo, nil)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "BA123", got[0].Number)
	})

	t.Run("sends dateLocalRole when provided", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.NotEmpty(t, r.URL.Query().Get("dateLocalRole"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
		}))
		defer srv.Close()

		dir := model.FlightDirection(model.FlightDirectionArrival)
		_, err := newTestClient(srv.URL).GetFlightScheduleHistory("Number", "BA123", dateFrom, dateTo, &dir)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"invalid parameters"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightScheduleHistory("Number", "BA123", dateFrom, dateTo, nil)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})
}

// ── GetFlightDepartureDates ───────────────────────────────────────────────────

func TestGetFlightDepartureDates(t *testing.T) {
	t.Run("returns date list with correct path", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/Number/BA123/dates", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]string{"2024-01-15", "2024-01-16"})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightDepartureDates("Number", "BA123")
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "2024-01-15", got[0])
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"flight not found"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightDepartureDates("Number", "XX999")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	})
}

// ── GetFlightDepartureDateRange ───────────────────────────────────────────────

func TestGetFlightDepartureDateRange(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 23, 59, 0, 0, time.UTC)

	t.Run("returns dates with correct path", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/Number/BA123/dates/2024-01-01T00:00/2024-01-31T23:59", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]string{"2024-01-05", "2024-01-12", "2024-01-19"})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightDepartureDateRange("Number", "BA123", from, to)
		require.NoError(t, err)
		assert.Len(t, got, 3)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"rate limit exceeded"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightDepartureDateRange("Number", "BA123", from, to)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
	})
}

// ── GetFlightInformationDepartureByLocalTime ──────────────────────────────────

func TestGetFlightInformationDepartureByLocalTime(t *testing.T) {
	fromLocal := time.Date(2024, 1, 15, 6, 0, 0, 0, time.UTC)
	toLocal := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	t.Run("sends correct path and optional params", func(t *testing.T) {
		boolTrue := true
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/airports/icao/EGLL/2024-01-15T06:00/2024-01-15T10:00", r.URL.Path)
			assert.Equal(t, "true", r.URL.Query().Get("withLeg"))
			assert.Equal(t, "true", r.URL.Query().Get("withCancelled"))
			w.Header().Set("Content-Type", "application/json")
			result := model.AirportFidsContract{
				Departures: []model.AirportFlightContract{{Number: "BA456"}},
				Arrivals:   []model.AirportFlightContract{},
			}
			_ = json.NewEncoder(w).Encode(result)
		}))
		defer srv.Close()

		dir := model.FlightDirection(model.FlightDirectionBoth)
		got, err := newTestClient(srv.URL).GetFlightInformationDepartureByLocalTime(
			"icao", "EGLL", fromLocal, toLocal,
			&dir, &boolTrue, &boolTrue, nil, nil, nil, nil,
		)
		require.NoError(t, err)
		require.Len(t, got.Departures, 1)
		assert.Equal(t, "BA456", got.Departures[0].Number)
	})

	t.Run("omits optional params when nil", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.Query().Get("withLeg"))
			assert.Empty(t, r.URL.Query().Get("direction"))
			assert.Empty(t, r.URL.Query().Get("offsetMinutes"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"departures":[],"arrivals":[]}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightInformationDepartureByLocalTime(
			"icao", "EGLL", fromLocal, toLocal,
			nil, nil, nil, nil, nil, nil, nil,
		)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"forbidden"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightInformationDepartureByLocalTime(
			"icao", "EGLL", fromLocal, toLocal,
			nil, nil, nil, nil, nil, nil, nil,
		)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	})
}

// ── GetFlightInformationDepartureByRelativeTime ───────────────────────────────

func TestGetFlightInformationDepartureByRelativeTime(t *testing.T) {
	t.Run("sends correct path and relative time params", func(t *testing.T) {
		offset := int64(30)
		duration := int64(120)
		boolFalse := false
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/airports/iata/LHR", r.URL.Path)
			assert.Equal(t, "30", r.URL.Query().Get("offsetMinutes"))
			assert.Equal(t, "120", r.URL.Query().Get("durationMinutes"))
			assert.Equal(t, "false", r.URL.Query().Get("withCargo"))
			w.Header().Set("Content-Type", "application/json")
			result := model.AirportFidsContract{
				Departures: []model.AirportFlightContract{{Number: "LH400"}},
				Arrivals:   []model.AirportFlightContract{{Number: "LH401"}},
			}
			_ = json.NewEncoder(w).Encode(result)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFlightInformationDepartureByRelativeTime(
			"iata", "LHR",
			nil, nil, nil, nil, &boolFalse, nil, nil,
			&offset, &duration,
		)
		require.NoError(t, err)
		require.Len(t, got.Departures, 1)
		require.Len(t, got.Arrivals, 1)
		assert.Equal(t, "LH400", got.Departures[0].Number)
	})

	t.Run("omits offset and duration when nil", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.Query().Get("offsetMinutes"))
			assert.Empty(t, r.URL.Query().Get("durationMinutes"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"departures":[],"arrivals":[]}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightInformationDepartureByRelativeTime(
			"iata", "LHR",
			nil, nil, nil, nil, nil, nil, nil,
			nil, nil,
		)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFlightInformationDepartureByRelativeTime(
			"iata", "LHR",
			nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
	})
}

// ── SearchFlightByTerm ────────────────────────────────────────────────────────

func TestSearchFlightByTerm(t *testing.T) {
	t.Run("sends term and limit as query params", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/flights/search/term", r.URL.Path)
			assert.Equal(t, "BA1", r.URL.Query().Get("q"))
			assert.Equal(t, "5", r.URL.Query().Get("limit"))
			w.Header().Set("Content-Type", "application/json")
			result := model.StringFlightSearchItemContractSearchResultCollectionContract{
				SearchBy: "Number",
				Count:    2,
				Items: []model.FlightSearchItemContract{
					{Number: "BA100"},
					{Number: "BA101"},
				},
			}
			_ = json.NewEncoder(w).Encode(result)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).SearchFlightByTerm("BA1", 5)
		require.NoError(t, err)
		assert.Equal(t, int64(2), got.Count)
		assert.Equal(t, "BA100", got.Items[0].Number)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"term too short"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).SearchFlightByTerm("B", 10)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "term too short", apiErr.Message)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{bad json`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).SearchFlightByTerm("BA1", 5)
		require.Error(t, err)
	})
}
