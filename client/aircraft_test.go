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

// newTestClient creates a Client pointed at the given test server URL.
func newTestClient(serverURL string) *Client {
	cfg := &model.Config{Host: "api.aerodatabox.com", Key: "test-key"}
	c, _ := NewClient(zap.NewNop(), cfg)
	c.URL = serverURL
	return c
}

// ── EndpointAircraft ─────────────────────────────────────────────────────────

func TestEndpointAircraft_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointAircraft
		want     string
	}{
		{EndpointAircraftSearchSingle, "/aircrafts/{searchBy}/{searchParameter}"},              //nolint:misspell
		{EndpointAircraftSearchAll, "/aircrafts/{searchBy}/{searchParameter}/all"},             //nolint:misspell
		{EndpointAircraftRegoHistory, "/aircrafts/{searchBy}/{searchParameter}/registrations"}, //nolint:misspell
		{EndpointAircraftSearchByTerm, "/aircrafts/search/term"},                               //nolint:misspell
		{EndpointAirlineFleetByCode, "/airlines/{code}/aircrafts"},                             //nolint:misspell
		{EndpointAircraftImageByTailNumber, "/aircrafts/reg/{reg}/image/beta"},                 //nolint:misspell
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String(), "endpoint %d", tc.endpoint)
	}
}

func TestEndpointAircraft_Replace(t *testing.T) {
	t.Run("all tokens substituted", func(t *testing.T) {
		got, err := EndpointAircraftSearchSingle.Replace(map[string]string{
			"searchBy":        "reg",
			"searchParameter": "G-ABCD",
		})
		require.NoError(t, err)
		assert.Equal(t, "/aircrafts/reg/G-ABCD", got) //nolint:misspell
	})

	t.Run("missing one token returns error naming it", func(t *testing.T) {
		_, err := EndpointAircraftSearchSingle.Replace(map[string]string{
			"searchBy": "reg",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "searchParameter")
	})

	t.Run("missing all tokens lists each one", func(t *testing.T) {
		_, err := EndpointAircraftSearchSingle.Replace(map[string]string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "searchBy")
		assert.Contains(t, err.Error(), "searchParameter")
	})

	t.Run("no-token endpoint accepts empty map", func(t *testing.T) {
		got, err := EndpointAircraftSearchByTerm.Replace(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "/aircrafts/search/term", got) //nolint:misspell
	})

	t.Run("extra keys in map are ignored", func(t *testing.T) {
		got, err := EndpointAirlineFleetByCode.Replace(map[string]string{
			"code":  "BA",
			"extra": "ignored",
		})
		require.NoError(t, err)
		assert.Equal(t, "/airlines/BA/aircrafts", got) //nolint:misspell
	})
}

// ── GetSingleAircraft ─────────────────────────────────────────────────────────

func TestGetSingleAircraft_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"aircraft not found"}`))
	}))
	defer server.Close()

	c := newTestClient(server.URL)
	result, err := c.GetSingleAircraft("reg", "XX-NONE", nil, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	var apiErr *model.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "aircraft not found", apiErr.Message)
}

func TestGetSingleAircraft_QueryParams(t *testing.T) {
	boolTrue := true

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "true", r.URL.Query().Get("withRegistrations"))
		assert.Equal(t, "true", r.URL.Query().Get("withImage"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"reg":"G-ABCD"}`))
	}))
	defer server.Close()

	c := newTestClient(server.URL)
	_, err := c.GetSingleAircraft("reg", "G-ABCD", &boolTrue, &boolTrue)
	require.NoError(t, err)
}

// ── GetAircraft ───────────────────────────────────────────────────────────────

func TestGetAircraft(t *testing.T) {
	boolTrue := true

	t.Run("returns list of aircraft and checks path", func(t *testing.T) {
		expected := []*model.AircraftContract{
			{Registration: "G-ABCD"},
			{Registration: "N12345"},
		}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/aircrafts/reg/G-ABCD/all", r.URL.Path) //nolint:misspell
			assert.Equal(t, "true", r.URL.Query().Get("withRegistrations"))
			assert.Empty(t, r.URL.Query().Get("withImage"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraft("reg", "G-ABCD", &boolTrue, nil)

		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, "G-ABCD", result[0].Registration)
		assert.Equal(t, "N12345", result[1].Registration)
	})

	t.Run("returns empty list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("[]"))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraft("reg", "NONE", nil, nil)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"invalid API key"}`))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraft("reg", "G-ABCD", nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraft("reg", "G-ABCD", nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

// ── GetAircraftHistory ────────────────────────────────────────────────────────

func TestGetAircraftHistory(t *testing.T) {
	t.Run("returns history with registrations", func(t *testing.T) {
		expected := &model.AircraftHistory{
			NumRegistrations: 2,
			Items: []model.AircraftRegistration{
				{Registration: "G-ABCD", Active: true},
				{Registration: "G-OLD", Active: false},
			},
		}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/aircrafts/reg/G-ABCD/registrations", r.URL.Path) //nolint:misspell
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraftHistory("reg", "G-ABCD")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(2), result.NumRegistrations)
		require.Len(t, result.Items, 2)
		assert.Equal(t, "G-ABCD", result.Items[0].Registration)
		assert.True(t, result.Items[0].Active)
		assert.Equal(t, "G-OLD", result.Items[1].Registration)
		assert.False(t, result.Items[1].Active)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message":"forbidden"}`))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraftHistory("reg", "G-ABCD")

		require.Error(t, err)
		assert.Nil(t, result)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraftHistory("reg", "G-ABCD")

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

// ── GetAircraftByTerm ─────────────────────────────────────────────────────────

func TestGetAircraftByTerm(t *testing.T) {
	t.Run("returns matching aircraft without limit", func(t *testing.T) {
		expected := &model.StringAircraftContractSearchResultCollectionContract{
			SearchBy: "Boeing",
			Count:    2,
			Items: []model.AircraftContract{
				{Registration: "G-ABCD"},
				{Registration: "N12345"},
			},
		}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/aircrafts/search/term", r.URL.Path) //nolint:misspell
			assert.Equal(t, "Boeing", r.URL.Query().Get("q"))
			assert.Empty(t, r.URL.Query().Get("limit"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraftByTerm("Boeing", nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(2), result.Count)
		assert.Len(t, result.Items, 2)
	})

	t.Run("sends limit query param when provided", func(t *testing.T) {
		limit := int64(5)
		expected := &model.StringAircraftContractSearchResultCollectionContract{}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "5", r.URL.Query().Get("limit"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		_, err := c.GetAircraftByTerm("Boeing", &limit)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"message":"rate limit exceeded"}`))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAircraftByTerm("Boeing", nil)

		require.Error(t, err)
		assert.Nil(t, result)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
		assert.Equal(t, "rate limit exceeded", apiErr.Message)
	})
}

// ── GetAirlineFleetByCode ─────────────────────────────────────────────────────

func TestGetAirlineFleetByCode(t *testing.T) {
	t.Run("returns paged fleet with all params", func(t *testing.T) {
		pageSize := int64(10)
		pageOffset := int64(20)
		withRegos := true

		expected := &model.AircraftContractPagedCollectionContract{
			TotalCount:  25,
			PageSize:    10,
			PageOffset:  20,
			HasNextPage: false,
			Items:       []model.AircraftContract{{Registration: "G-ABCD"}},
		}
		jsonData, _ := json.Marshal(expected)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/airlines/BA/aircrafts", r.URL.Path) //nolint:misspell
			assert.Equal(t, "10", r.URL.Query().Get("pageSize"))
			assert.Equal(t, "20", r.URL.Query().Get("pageOffset"))
			assert.Equal(t, "true", r.URL.Query().Get("withRegistrations"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAirlineFleetByCode("BA", &pageSize, &pageOffset, &withRegos)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(25), result.TotalCount)
		assert.False(t, result.HasNextPage)
		require.Len(t, result.Items, 1)
		assert.Equal(t, "G-ABCD", result.Items[0].Registration)
	})

	t.Run("omits optional params when nil", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.URL.Query().Get("pageSize"))
			assert.Empty(t, r.URL.Query().Get("pageOffset"))
			assert.Empty(t, r.URL.Query().Get("withRegistrations"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"totalCount":0,"pageOffset":0,"pageSize":0,"hasNextPage":false,"items":[]}`))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAirlineFleetByCode("BA", nil, nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Items)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"airline not found"}`))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAirlineFleetByCode("XX", nil, nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		c := newTestClient(server.URL)
		result, err := c.GetAirlineFleetByCode("BA", nil, nil, nil)

		require.Error(t, err)
		assert.Nil(t, result)
	})
}
