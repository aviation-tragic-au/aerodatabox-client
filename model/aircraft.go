package model

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
	"time"
)

// AircraftParams represents the parameters for the aircraft API
type AircraftParams struct {
	WithRegistrations *bool `url:"withRegistrations,omitempty"`
	WithImage         *bool `url:"withImage,omitempty"`
}

// NewAircraftParams creates a new aircraft params instance
func NewAircraftParams(withRegistrations, withImage *bool) *AircraftParams {
	return &AircraftParams{
		WithRegistrations: withRegistrations,
		WithImage:         withImage,
	}
}

// AirlineFleetByCodeParams represents the parameters for the airline fleet by code API
type AirlineFleetByCodeParams struct {
	PageSize          *int64 `url:"pageSize,omitempty"`
	PageOffset        *int64 `url:"pageOffset,omitempty"`
	WithRegistrations *bool  `url:"withRegistrations,omitempty"`
}

// NewAirlineFleetByCodeParams creates a new airline fleet by code params instance
func NewAirlineFleetByCodeParams(pageSize, pageOffset *int64, withRegistrations *bool) *AirlineFleetByCodeParams {
	return &AirlineFleetByCodeParams{
		PageSize:          pageSize,
		PageOffset:        pageOffset,
		WithRegistrations: withRegistrations,
	}
}

// AircraftSearchParams represents the parameters for the aircraft search by term API
type AircraftSearchParams struct {
	Term  string `url:"q"`
	Limit *int64 `url:"limit,omitempty"`
}

// NewAircraftSearchParams creates a new aircraft search params instance
func NewAircraftSearchParams(term string, limit *int64) *AircraftSearchParams {
	return &AircraftSearchParams{
		Term:  term,
		Limit: limit,
	}
}

// AircraftContract represents an aircraft
type AircraftContract struct {
	ID               int64                  `json:"id"`
	Registration     string                 `json:"reg"`
	Active           bool                   `json:"active"`
	Serial           *string                `json:"serial"`
	HexICAO          *string                `json:"hexIcao"`
	AirlineName      *string                `json:"airlineName"`
	IATAType         *string                `json:"iataType"`
	IATACodeShort    *string                `json:"iataCodeShort"`
	ICAOCode         *string                `json:"icaoCode"`
	Model            *string                `json:"model"`
	ModelCode        *string                `json:"modelCode"`
	NumSeats         *int64                 `json:"numSeats"`
	RolloutDate      *time.Time             `json:"rolloutDate"`
	FirstFlightDate  *time.Time             `json:"firstFlightDate"`
	DeliveryDate     *time.Time             `json:"deliveryDate"`
	RegistrationDate *time.Time             `json:"registrationDate"`
	TypeName         *string                `json:"typeName"`
	NumEngines       *int64                 `json:"numEngines"`
	EngineType       string                 `json:"engineType"`
	IsFreighter      bool                   `json:"isFreighter"`
	ProductionLine   *string                `json:"productionLine"`
	AgeYears         *float64               `json:"ageYears"`
	Verified         bool                   `json:"verified"`
	Image            *ResourceContract      `json:"image"`
	NumRegistrations int64                  `json:"numRegistrations"`
	Registrations    []AircraftRegistration `json:"registrations"`
}

// ResourceContract represents an aircraft image
type ResourceContract struct {
	URL              string      `json:"url"`
	WebURL           *string     `json:"webUrl"`
	Author           *string     `json:"author"`
	Title            *string     `json:"title"`
	Description      *string     `json:"description"`
	License          LicenseType `json:"license"`
	HTMLAttributions []*string   `json:"htmlAttributions"`
}

// AircraftRegistration represents an aircraft registration
type AircraftRegistration struct {
	Registration     string     `json:"reg"`
	Active           bool       `json:"active"`
	HexICAO          *string    `json:"hexIcao"`
	AirlineName      *string    `json:"airlineName"`
	RegistrationDate *time.Time `json:"registrationDate"`
}

// AircraftPaged represents a paged aircraft
type AircraftPaged struct {
	TotalCount  int64              `json:"totalCount"`
	PageOffset  int64              `json:"pageOffset"`
	PageSize    int64              `json:"pageSize"`
	HasNextPage bool               `json:"hasNextPage"`
	Count       int64              `json:"count"`
	Items       []AircraftContract `json:"items"`
}

// Airframe represents a simplified aircraft record returned by search endpoints.
type Airframe struct {
	Registration  string  `json:"reg"`
	HexICAO       string  `json:"hexIcao"`
	TypeName      string  `json:"typeName"`
	IATAType      *string `json:"iataType"`
	IATACodeShort *string `json:"iataCodeShort"`
	ICAOCode      *string `json:"icaoCode"`
	Model         *string `json:"model"`
	ModelCode     *string `json:"modelCode"`
}

// AircraftHistory represents the registration history of an aircraft.
type AircraftHistory struct {
	NumRegistrations int64                  `json:"numRegistrations"`
	Items            []AircraftRegistration `json:"registrations"`
}

// StringAircraftContractSearchResultCollectionContract represents a collection of aircraft contracts.
type StringAircraftContractSearchResultCollectionContract struct {
	SearchBy string             `json:"searchBy"`
	Count    int64              `json:"count"`
	Items    []AircraftContract `json:"items"`
}

// AircraftContractPagedCollectionContract represents a paged collection of aircraft contracts.
type AircraftContractPagedCollectionContract struct {
	TotalCount  int64              `json:"totalCount"`
	PageOffset  int64              `json:"pageOffset"`
	PageSize    int64              `json:"pageSize"`
	HasNextPage bool               `json:"hasNextPage"`
	Items       []AircraftContract `json:"items"`
}
