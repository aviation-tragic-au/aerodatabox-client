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

// ── EndpointIndustry ──────────────────────────────────────────────────────────

func TestEndpointIndustry_String(t *testing.T) {
	assert.Equal(t, "/industry/faa-ladd/{id}/status", EndpointIndustryFAALADDAircraftStatus.String())
}

func TestEndpointIndustry_Replace(t *testing.T) {
	t.Run("substitutes id token", func(t *testing.T) {
		got, err := EndpointIndustryFAALADDAircraftStatus.Replace(map[string]string{
			"id": "N12345",
		})
		require.NoError(t, err)
		assert.Equal(t, "/industry/faa-ladd/N12345/status", got)
	})

	t.Run("missing id token returns error", func(t *testing.T) {
		_, err := EndpointIndustryFAALADDAircraftStatus.Replace(map[string]string{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "id")
	})
}

// ── GetFAALADDAircraftStatus ──────────────────────────────────────────────────

func TestGetFAALADDAircraftStatus(t *testing.T) {
	t.Run("returns status for unblocked aircraft", func(t *testing.T) {
		expected := &model.FAALADDAircraftStatusContract{
			ID:        "N12345",
			IsBlocked: false,
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/industry/faa-ladd/N12345/status", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFAALADDAircraftStatus("N12345")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "N12345", got.ID)
		assert.False(t, got.IsBlocked)
		assert.Nil(t, got.BlockedSince)
	})

	t.Run("returns status for blocked aircraft with timestamps", func(t *testing.T) {
		blockedAt := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
		expected := &model.FAALADDAircraftStatusContract{
			ID:            "N99999",
			IsBlocked:     true,
			BlockedSince:  &blockedAt,
			LastBlockedOn: &blockedAt,
		}
		jsonData, _ := json.Marshal(expected)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/industry/faa-ladd/N99999/status", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonData)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetFAALADDAircraftStatus("N99999")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.True(t, got.IsBlocked)
		require.NotNil(t, got.BlockedSince)
		assert.Equal(t, blockedAt, *got.BlockedSince)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetFAALADDAircraftStatus("UNKNOWN")
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

		_, err := newTestClient(srv.URL).GetFAALADDAircraftStatus("N12345")
		require.Error(t, err)
	})
}
