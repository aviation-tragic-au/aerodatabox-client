package model

import "time"

// FlightSubscriberContract represents the subscriber contract for the flight alert API
type FlightSubscriberContract struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// FlightSubjectContract represents the subject contract for the flight alert API
type FlightSubjectContract struct {
	Type FlightSubscriptionSubjectType `json:"type"`
	ID   string                        `json:"id"`
}

// FlightSubscriptionBalanceContract represents the subscription balance contract for the flight alert API
type FlightSubscriptionBalanceContract struct {
	CreditsRemaining int64      `json:"creditsRemaining"`
	LastRefilledUTC  *time.Time `json:"lastRefilledUtc"`
	LastDeductedUTC  *time.Time `json:"lastDeductedUtc"`
}

// FlightSubscriptionContract represents the subscription contract for the flight alert API
type FlightSubscriptionContract struct {
	ID                string                        `json:"id"`
	IsActive          bool                          `json:"isActive"`
	BillingType       FlightSubscriptionBillingType `json:"billingType"`
	ActivateBeforeUTC *time.Time                    `json:"activateBeforeUtc"`
	ExpiresOnUTC      *time.Time                    `json:"expiresOnUtc"`
	CreatedOnUTC      *time.Time                    `json:"createdOnUtc"`
	Subject           FlightSubjectContract         `json:"subject"`
	Subscriber        FlightSubscriberContract      `json:"subscriber"`
	Notices           []string                      `json:"notices"`
}

// FlightSubscriptionBody represents the body of a flight alert subscription
type FlightSubscriptionBody struct {
	URL                string `json:"url"`
	MaxDeliveryRetries *int64 `json:"maxDeliveryRetries,omitempty"`
}

// FlightSubscriptionRefillBody represents the body of a flight alert subscription refill
type FlightSubscriptionRefillBody struct {
	Credits int64 `json:"credits"`
}
