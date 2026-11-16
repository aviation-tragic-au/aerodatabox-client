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

	"github.com/aviation-tragic-au/aerodatabox-client/model"
)

// EndpointFlightAlert represents the flight alert API endpoints
type EndpointFlightAlert int

const (
	EndpointCreateWebhookSubscription EndpointFlightAlert = iota
	EndpointGetWebhookSubscription
	EndpointDeleteWebhookSubscription
	EndpointListWebhookSubscriptions
	EndpointGetWebhookSubscriptionBalance
	EndpointRefillWebhookSubscriptionBalance
)

var endpointFlightAlert = map[EndpointFlightAlert]string{
	EndpointCreateWebhookSubscription:        "/subscriptions/webhook/{subjectType}/{subjectId}",
	EndpointGetWebhookSubscription:           "/subscriptions/webhook/{subscriptionId}",
	EndpointDeleteWebhookSubscription:        "/subscriptions/webhook/{subscriptionId}",
	EndpointListWebhookSubscriptions:         "/subscriptions/webhook",
	EndpointGetWebhookSubscriptionBalance:    "/subscriptions/balance",
	EndpointRefillWebhookSubscriptionBalance: "/subscriptions/balance/refill",
}

// String returns the endpoint URL.
func (e EndpointFlightAlert) String() string {
	return endpointFlightAlert[e]
}

// Replace returns the endpoint with tokens replaced.
func (e EndpointFlightAlert) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// CreateWebhookSubscription creates a new webhook subscription
func (c *Client) CreateWebhookSubscription(ctx context.Context, subjectID string, subjectType model.FlightSubscriptionSubjectType, callbackURL string, maxDeliveryRetries *int64) (*model.FlightSubscriberContract, error) {
	endpoint, err := EndpointCreateWebhookSubscription.Replace(map[string]string{
		"subjectId":   subjectID,
		"subjectType": subjectType.String(),
	})
	if err != nil {
		return nil, err
	}

	req, err := c.postRequest(ctx, endpoint, &model.FlightSubscriptionBody{
		URL:                callbackURL,
		MaxDeliveryRetries: maxDeliveryRetries,
	})
	if err != nil {
		return nil, err
	}

	var output *model.FlightSubscriberContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetWebhookSubscription retrieves a webhook subscription by its ID
func (c *Client) GetWebhookSubscription(ctx context.Context, subscriptionID string) (*model.FlightSubscriberContract, error) {
	url, err := EndpointGetWebhookSubscription.Replace(map[string]string{
		"subscriptionId": subscriptionID,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var output *model.FlightSubscriberContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// DeleteWebhookSubscription deletes a webhook subscription by its ID
func (c *Client) DeleteWebhookSubscription(ctx context.Context, subscriptionID string) error {
	url, err := EndpointDeleteWebhookSubscription.Replace(map[string]string{
		"subscriptionId": subscriptionID,
	})
	if err != nil {
		return err
	}

	req, err := c.deleteRequest(ctx, url)
	if err != nil {
		return err
	}

	return c.executeRequest(req, nil)
}

// ListWebhookSubscriptions lists all webhook subscriptions
func (c *Client) ListWebhookSubscriptions(ctx context.Context) ([]*model.FlightSubscriberContract, error) {
	url := EndpointListWebhookSubscriptions.String()
	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var output []*model.FlightSubscriberContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetWebhookSubscriptionsBalance retrieves the webhook subscription balance
func (c *Client) GetWebhookSubscriptionsBalance(ctx context.Context) (*model.FlightSubscriptionBalanceContract, error) {
	url := EndpointGetWebhookSubscriptionBalance.String()
	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	var output *model.FlightSubscriptionBalanceContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// RefillWebhookSubscriptionsBalance refills the webhook subscription balance
func (c *Client) RefillWebhookSubscriptionsBalance(ctx context.Context, credits int64) (*model.FlightSubscriptionBalanceContract, error) {
	url := EndpointRefillWebhookSubscriptionBalance.String()
	req, err := c.postRequest(ctx, url, &model.FlightSubscriptionRefillBody{Credits: credits})
	if err != nil {
		return nil, err
	}
	var output *model.FlightSubscriptionBalanceContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}
