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

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── EndpointAirport ───────────────────────────────────────────────────────────

func TestEndpointAirport_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointAirport
		want     string
	}{
		{EndpointAirportByCode, "/airports/{codeType}/{code}"},
		{EndpointAirportRunways, "/airports/{codeType}/{code}/runways"},
		{EndpointAirportSearchByLocation, "/airports/search/location"},
		{EndpointAirportSearchByIP, "/airports/search/ip"},
		{EndpointAirportSearchByTerm, "/airports/search/term"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String())
	}
}

func TestEndpointAirport_Replace(t *testing.T) {
	t.Run("all tokens substituted", func(t *testing.T) {
		got, err := EndpointAirportByCode.Replace(map[string]string{
			"codeType": "icao",
			"code":     "EGLL",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/icao/EGLL", got)
	})

	t.Run("missing token returns error", func(t *testing.T) {
		_, err := EndpointAirportByCode.Replace(map[string]string{"codeType": "icao"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "code")
	})

	t.Run("no-token endpoint accepts empty map", func(t *testing.T) {
		got, err := EndpointAirportSearchByLocation.Replace(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "/airports/search/location", got)
	})

	t.Run("extra keys are ignored", func(t *testing.T) {
		got, err := EndpointAirportByCode.Replace(map[string]string{
			"codeType": "iata",
			"code":     "LHR",
			"extra":    "ignored",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airports/iata/LHR", got)
	})
}

// ── GetAirportByCode ──────────────────────────────────────────────────────────

func TestGetAirportByCode_QueryParams(t *testing.T) {
	boolTrue := true

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/airports/icao/EGLL", r.URL.Path)
		assert.Equal(t, "true", r.URL.Query().Get("withTime"))
		assert.Equal(t, "true", r.URL.Query().Get("withRunways"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"Heathrow"}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	_, err := c.GetAirportByCode(context.Background(), "EGLL", "icao", &boolTrue, &boolTrue)
	require.NoError(t, err)
}

func TestGetAirportByCode(t *testing.T) {
	t.Run("returns airport and checks path", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/iata/LHR", r.URL.Path)
			assert.Empty(t, r.URL.Query().Get("withTime"))
			assert.Empty(t, r.URL.Query().Get("withRunways"))
			w.Header().Set("Content-Type", "application/json")
			airport := model.Airport{Name: "Heathrow", CountryCode: "GB"}
			_ = json.NewEncoder(w).Encode(airport)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportByCode(context.Background(), "LHR", "iata", nil, nil)
		require.NoError(t, err)
		assert.Equal(t, "Heathrow", got.Name)
		assert.Equal(t, "GB", got.CountryCode)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"airport not found"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportByCode(context.Background(), "ZZZ", "icao", nil, nil)
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

		_, err := newTestClient(srv.URL).GetAirportByCode(context.Background(), "LHR", "iata", nil, nil)
		require.Error(t, err)
	})
}

// ── GetAirportRunways ─────────────────────────────────────────────────────────

func TestGetAirportRunways(t *testing.T) {
	t.Run("returns runway list", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/icao/EGLL/runways", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			runways := []*model.RunwayContract{{Name: "09L/27R"}, {Name: "09R/27L"}}
			_ = json.NewEncoder(w).Encode(runways)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportRunways(context.Background(), "EGLL", "icao")
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "09L/27R", got[0].Name)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"invalid code"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportRunways(context.Background(), "BAD", "icao")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})
}

// ── GetAirportByLocation ──────────────────────────────────────────────────────

func TestGetAirportByLocation(t *testing.T) {
	t.Run("sends all query params", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/search/location", r.URL.Path)
			assert.Equal(t, "51.477", r.URL.Query().Get("lat"))
			assert.Equal(t, "-0.461", r.URL.Query().Get("lon"))
			assert.Equal(t, "5", r.URL.Query().Get("limit"))
			assert.Equal(t, "50", r.URL.Query().Get("radiusKm"))
			assert.Equal(t, "true", r.URL.Query().Get("withFlightInfoOnly"))
			w.Header().Set("Content-Type", "application/json")
			result := model.StringListingAirportContractSearchResultCollectionContract{Count: 1}
			_ = json.NewEncoder(w).Encode(result)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportByLocation(context.Background(), 51.477, -0.461, 5, 50, true)
		require.NoError(t, err)
		assert.Equal(t, int64(1), got.Count)
	})

	t.Run("omits withFlightInfoOnly when false", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.Query().Get("withFlightInfoOnly"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"count":0,"items":[]}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportByLocation(context.Background(), 51.477, -0.461, 5, 50, false)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportByLocation(context.Background(), 0, 0, 5, 50, false)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})
}

// ── GetAirportByIP ────────────────────────────────────────────────────────────

func TestGetAirportByIP(t *testing.T) {
	t.Run("sends all query params", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/search/ip", r.URL.Path)
			assert.Equal(t, "8.8.8.8", r.URL.Query().Get("q"))
			assert.Equal(t, "10", r.URL.Query().Get("limit"))
			assert.Equal(t, "100", r.URL.Query().Get("radiusKm"))
			assert.Equal(t, "true", r.URL.Query().Get("withFlightInfoOnly"))
			w.Header().Set("Content-Type", "application/json")
			result := model.StringListingAirportContractSearchResultCollectionContract{Count: 2}
			_ = json.NewEncoder(w).Encode(result)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportByIP(context.Background(), "8.8.8.8", 10, 100, true)
		require.NoError(t, err)
		assert.Equal(t, int64(2), got.Count)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"invalid IP"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportByIP(context.Background(), "not-an-ip", 5, 50, false)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})
}

// ── GetAirportByTerm ──────────────────────────────────────────────────────────

func TestGetAirportByTerm(t *testing.T) {
	t.Run("sends term without optional params", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airports/search/term", r.URL.Path)
			assert.Equal(t, "heathrow", r.URL.Query().Get("q"))
			assert.Empty(t, r.URL.Query().Get("limit"))
			assert.Empty(t, r.URL.Query().Get("withFlightInfoOnly"))
			assert.Empty(t, r.URL.Query().Get("withSearchByCode"))
			w.Header().Set("Content-Type", "application/json")
			result := model.StringListingAirportContractSearchResultCollectionContract{
				SearchBy: "term",
				Count:    1,
				Items:    []model.ListingAirportContract{{Name: "London Heathrow"}},
			}
			_ = json.NewEncoder(w).Encode(result)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetAirportByTerm(context.Background(), "heathrow", nil, false, false)
		require.NoError(t, err)
		assert.Equal(t, int64(1), got.Count)
		assert.Equal(t, "London Heathrow", got.Items[0].Name)
	})

	t.Run("sends limit and boolean flags when set", func(t *testing.T) {
		limit := int64(3)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "3", r.URL.Query().Get("limit"))
			assert.Equal(t, "true", r.URL.Query().Get("withFlightInfoOnly"))
			assert.Equal(t, "true", r.URL.Query().Get("withSearchByCode"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"count":0,"items":[]}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportByTerm(context.Background(), "LHR", &limit, true, true)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"rate limit exceeded"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetAirportByTerm(context.Background(), "london", nil, false, false)
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

		_, err := newTestClient(srv.URL).GetAirportByTerm(context.Background(), "london", nil, false, false)
		require.Error(t, err)
	})
}
