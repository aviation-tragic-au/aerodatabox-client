package model

import "time"

// AirportFeedServiceStatusContract represents the status of the airport feed services
type AirportFeedServiceStatusContract struct {
	FlightSchedulesFeed FeedServiceStatusContract `json:"flightSchedulesFeed"`
	FlightLiveUpdates   FeedServiceStatusContract `json:"flightLiveUpdates"`
	ADSBUpdates         FeedServiceStatusContract `json:"adsbUpdates"`
	GeneralAvailability GeneralAvailability       `json:"generalAvailability"`
}

// FeedServiceStatusContract represents the status of a feed service
type FeedServiceStatusContract struct {
	Service               FeedServiceEnum   `json:"service"`
	Status                FeedServiceStatus `json:"status"`
	MinAvailableLocalDate *time.Time        `json:"minAvailableLocalDate"`
	MaxAvailableLocalDate *time.Time        `json:"maxAvailableLocalDate"`
}

// GeneralAvailability represents the general availability of a feed service
type GeneralAvailability struct {
	MinAvailableLocalDate *time.Time `json:"minAvailableLocalDate"`
	MaxAvailableLocalDate *time.Time `json:"maxAvailableLocalDate"`
}

// StringCollectionContract represents a collection of strings
type StringCollectionContract struct {
	Count int64    `json:"count"`
	Items []string `json:"items"`
}
