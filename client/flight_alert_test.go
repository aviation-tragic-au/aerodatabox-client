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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── EndpointFlightAlert ───────────────────────────────────────────────────────

func TestEndpointFlightAlert_String(t *testing.T) {
	cases := []struct {
		endpoint EndpointFlightAlert
		want     string
	}{
		{EndpointCreateWebhookSubscription, "/subscriptions/webhook/{subjectType}/{subjectId}"},
		{EndpointGetWebhookSubscription, "/subscriptions/webhook/{subscriptionId}"},
		{EndpointDeleteWebhookSubscription, "/subscriptions/webhook/{subscriptionId}"},
		{EndpointListWebhookSubscriptions, "/subscriptions/webhook"},
		{EndpointGetWebhookSubscriptionBalance, "/subscriptions/balance"},
		{EndpointRefillWebhookSubscriptionBalance, "/subscriptions/balance/refill"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.endpoint.String())
	}
}

func TestEndpointFlightAlert_Replace(t *testing.T) {
	t.Run("substitutes subjectType and subjectId", func(t *testing.T) {
		got, err := EndpointCreateWebhookSubscription.Replace(map[string]string{
			"subjectType": "FlightByNumber",
			"subjectId":   "BA123",
		})
		require.NoError(t, err)
		assert.Equal(t, "/subscriptions/webhook/FlightByNumber/BA123", got)
	})

	t.Run("substitutes subscriptionId", func(t *testing.T) {
		got, err := EndpointGetWebhookSubscription.Replace(map[string]string{
			"subscriptionId": "sub-abc-123",
		})
		require.NoError(t, err)
		assert.Equal(t, "/subscriptions/webhook/sub-abc-123", got)
	})

	t.Run("missing token returns error naming it", func(t *testing.T) {
		_, err := EndpointCreateWebhookSubscription.Replace(map[string]string{
			"subjectType": "FlightByNumber",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "subjectId")
	})

	t.Run("no-token endpoint accepts empty map", func(t *testing.T) {
		got, err := EndpointListWebhookSubscriptions.Replace(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "/subscriptions/webhook", got)
	})
}

// ── CreateWebhookSubscription ─────────────────────────────────────────────────

func TestCreateWebhookSubscription(t *testing.T) {
	t.Run("POSTs to correct path with body", func(t *testing.T) {
		retries := int64(3)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/subscriptions/webhook/FlightByNumber/BA123", r.URL.Path)

			var body model.FlightSubscriptionBody
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, "https://example.com/hook", body.URL)
			require.NotNil(t, body.MaxDeliveryRetries)
			assert.Equal(t, int64(3), *body.MaxDeliveryRetries)

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model.FlightSubscriberContract{Type: "webhook", ID: "sub-1"})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).CreateWebhookSubscription(context.Background(),
			"BA123", model.FlightSubscriptionSubjectTypeFlightByNumber,
			"https://example.com/hook", &retries,
		)
		require.NoError(t, err)
		assert.Equal(t, "sub-1", got.ID)
		assert.Equal(t, "webhook", got.Type)
	})

	t.Run("nil maxDeliveryRetries omitted from body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body model.FlightSubscriptionBody
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Nil(t, body.MaxDeliveryRetries)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model.FlightSubscriberContract{ID: "sub-2"})
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).CreateWebhookSubscription(context.Background(),
			"EGLL", model.FlightSubscriptionSubjectTypeFlightByAirportICAO,
			"https://example.com/hook", nil,
		)
		require.NoError(t, err)
	})

	t.Run("uses subjectType String() in path", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/subscriptions/webhook/FlightByAirportIcao/EGLL", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model.FlightSubscriberContract{ID: "sub-3"})
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).CreateWebhookSubscription(context.Background(),
			"EGLL", model.FlightSubscriptionSubjectTypeFlightByAirportICAO,
			"https://example.com/hook", nil,
		)
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"message":"invalid subject"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).CreateWebhookSubscription(context.Background(),
			"BAD", model.FlightSubscriptionSubjectTypeFlightByNumber,
			"https://example.com/hook", nil,
		)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnprocessableEntity, apiErr.StatusCode)
	})
}

// ── GetWebhookSubscription ────────────────────────────────────────────────────

func TestGetWebhookSubscription(t *testing.T) {
	t.Run("GETs correct path and returns subscription", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/subscriptions/webhook/sub-abc", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model.FlightSubscriberContract{Type: "webhook", ID: "sub-abc"})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetWebhookSubscription(context.Background(), "sub-abc")
		require.NoError(t, err)
		assert.Equal(t, "sub-abc", got.ID)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"subscription not found"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetWebhookSubscription(context.Background(), "missing")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	})
}

// ── DeleteWebhookSubscription ─────────────────────────────────────────────────

func TestDeleteWebhookSubscription(t *testing.T) {
	t.Run("DELETEs correct path and returns nil on 204", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/subscriptions/webhook/sub-xyz", r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer srv.Close()

		err := newTestClient(srv.URL).DeleteWebhookSubscription(context.Background(), "sub-xyz")
		require.NoError(t, err)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"subscription not found"}`))
		}))
		defer srv.Close()

		err := newTestClient(srv.URL).DeleteWebhookSubscription(context.Background(), "missing")
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	})
}

// ── ListWebhookSubscriptions ──────────────────────────────────────────────────

func TestListWebhookSubscriptions(t *testing.T) {
	t.Run("returns list of subscriptions", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/subscriptions/webhook", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			subs := []*model.FlightSubscriberContract{
				{Type: "webhook", ID: "sub-1"},
				{Type: "webhook", ID: "sub-2"},
			}
			_ = json.NewEncoder(w).Encode(subs)
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).ListWebhookSubscriptions(context.Background())
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "sub-1", got[0].ID)
		assert.Equal(t, "sub-2", got[1].ID)
	})

	t.Run("returns empty list", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).ListWebhookSubscriptions(context.Background())
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).ListWebhookSubscriptions(context.Background())
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	})
}

// ── GetWebhookSubscriptionsBalance ───────────────────────────────────────────

func TestGetWebhookSubscriptionsBalance(t *testing.T) {
	t.Run("returns balance", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/subscriptions/balance", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model.FlightSubscriptionBalanceContract{CreditsRemaining: 42})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).GetWebhookSubscriptionsBalance(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(42), got.CreditsRemaining)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusPaymentRequired)
			_, _ = w.Write([]byte(`{"message":"no credits"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).GetWebhookSubscriptionsBalance(context.Background())
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusPaymentRequired, apiErr.StatusCode)
	})
}

// ── RefillWebhookSubscriptionsBalance ────────────────────────────────────────

func TestRefillWebhookSubscriptionsBalance(t *testing.T) {
	t.Run("POSTs credits and returns updated balance", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/subscriptions/balance/refill", r.URL.Path)

			body, _ := io.ReadAll(r.Body)
			assert.JSONEq(t, `{"credits":100}`, string(body))

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(model.FlightSubscriptionBalanceContract{CreditsRemaining: 142})
		}))
		defer srv.Close()

		got, err := newTestClient(srv.URL).RefillWebhookSubscriptionsBalance(context.Background(), 100)
		require.NoError(t, err)
		assert.Equal(t, int64(142), got.CreditsRemaining)
	})

	t.Run("API error bubbles up", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"message":"invalid credit amount"}`))
		}))
		defer srv.Close()

		_, err := newTestClient(srv.URL).RefillWebhookSubscriptionsBalance(context.Background(), -1)
		require.Error(t, err)
		var apiErr *model.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	})
}
