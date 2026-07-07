package model

import "time"

// This file contains models related to the Statistical API endpoints for Aerodatabox

// AirportDelayContract represents the airport delay contract
type AirportDelayContract struct {
	AirportICAO                string           `json:"airportIcao"`
	From                       Time             `json:"from"`
	To                         Time             `json:"to"`
	DeparturesDelayInformation DelayInformation `json:"departuresDelayInformation"`
	ArrivalsDelayInformation   DelayInformation `json:"arrivalsDelayInformation"`
}

// DelayInformation represents the delay information
type DelayInformation struct {
	Total          int64    `json:"numTotal"`
	QualifiedTotal int64    `json:"numQualifiedTotal"`
	Cancelled      int64    `json:"numCancelled"`
	MedianDelay    *string  `json:"medianDelay"`
	DelayIndex     *float64 `json:"delayIndex"`
}

// GlobalDelay represents the global delay information
type GlobalDelay struct {
	OngoingDelays []DelayInformation `json:"ongoingDelays"`
}

// DailyRouteStatContract represents the routes for an airport
type DailyRouteStatContract struct {
	Routes []DailyRouteStatRecordContract `json:"routes"`
}

// DailyRouteStatRecordContract represents a route
type DailyRouteStatRecordContract struct {
	Destination         AirportLocation         `json:"destination"`
	AverageDailyFlights float64                 `json:"averageDailyFlights"`
	Operations          []FlightAirlineContract `json:"operations"`
}

// FlightAirlineContract represents an airline
type FlightAirlineContract struct {
	Name string  `json:"name"`
	IATA *string `json:"iata"`
	ICAO *string `json:"icao"`
}

// FlightLegDelayContract represents a flight delay
type FlightLegDelayContract struct {
	Number       string                `json:"number"`
	Origins      []FlightDelayContract `json:"origins"`
	Destinations []FlightDelayContract `json:"destinations"`
}

// FlightDelayContract represents a flight delay airport
type FlightDelayContract struct {
	AirportICAO               string                      `json:"airportIcao"`
	Class                     string                      `json:"class"`
	ScheduledHourUTC          *int64                      `json:"scheduledHourUtc"`
	MedianDelay               string                      `json:"medianDelay"`
	DelayPercentiles          []PercentileBracketContract `json:"delayPercentiles"`
	NumConsideredFlights      int64                       `json:"numConsideredFlights"`
	NumFlightsDelayedBrackets []DelayBracketContract      `json:"numFlightsDelayedBrackets"`
	FromUTC                   time.Time                   `json:"fromUtc,omitempty"`
	ToUTC                     time.Time                   `json:"toUtc,omitempty"`
}

// PercentileBracketContract represents a delay percentile
type PercentileBracketContract struct {
	Percentile int64
	Delay      string
}

// DelayBracketContract represents a delay bracket
type DelayBracketContract struct {
	DelayedTo   string  `json:"delayedTo"`
	DelayedFrom *string `json:"delayedFrom,omitempty"`
	Number      int64   `json:"num"`
	Percentage  float64 `json:"percentage"`
}
