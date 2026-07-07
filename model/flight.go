package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type FlightAirportMovementQualityEnum int

const (
	FlightAirportMovementQualityEnumBasic FlightAirportMovementQualityEnum = iota
	FlightAirportMovementQualityEnumLive
	FlightAirportMovementQualityEnumApproximate
)

var flightAirportMovementQuality = map[FlightAirportMovementQualityEnum]string{
	FlightAirportMovementQualityEnumBasic:       "Basic",
	FlightAirportMovementQualityEnumLive:        "Live",
	FlightAirportMovementQualityEnumApproximate: "Approximate",
}

var flightAirportMovementQualityByString = map[string]FlightAirportMovementQualityEnum{
	"Basic":       FlightAirportMovementQualityEnumBasic,
	"Live":        FlightAirportMovementQualityEnumLive,
	"Approximate": FlightAirportMovementQualityEnumApproximate,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightAirportMovementQualityEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightAirportMovementQuality[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface for FlightAirportMovementQualityEnum
func (s *FlightAirportMovementQualityEnum) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	val, ok := flightAirportMovementQualityByString[str]
	if !ok {
		return fmt.Errorf("unknown FlightAirportMovementQualityEnum value: %q", str)
	}
	*s = val
	return nil
}

func (s *FlightAirportMovementQualityEnum) String() string {
	return flightAirportMovementQuality[*s]
}

// FlightContract represents a flight contract
type FlightContract struct {
	GreatCircleDistance *Measurement                  `json:"greatCircleDistance"`
	Departure           FlightAirportMovementContract `json:"departure"`
	Arrival             FlightAirportMovementContract `json:"arrival"`
	FlightPlan          FlightPlanContract            `json:"flightPlan"`
	LastUpdatedUTC      time.Time                     `json:"lastUpdatedUtc"`
	Number              string                        `json:"number"`
	CallSign            *string                       `json:"callSign"`
	Status              FlightStatus                  `json:"status"`
	CodeshareStatus     CodeshareStatus               `json:"codeshareStatus"`
	IsCargo             bool                          `json:"isCargo"`
	Aircraft            *FlightAircraftContract       `json:"aircraft"`
	Airline             *FlightAirlineContract        `json:"airline"`
	Location            *FlightLocationContract       `json:"location"`
}

// FlightAirportMovementContract represents the airport movement contract for a flight
type FlightAirportMovementContract struct {
	Airport       ListingAirportContract `json:"airport"`
	ScheduledTime *Time                  `json:"scheduledTime"`
	RevisedTime   *Time                  `json:"revisedTime"`
	PredictedTime *Time                  `json:"predictedTime"`
	RunwayTime    *Time                  `json:"runwayTime"`
	Terminal      *string                `json:"terminal"`
	CheckInDesk   *string                `json:"checkInDesk"`
	Gate          *string                `json:"gate"`
	BaggageBelt   *string                `json:"baggageBelt"`
	Runway        *string                `json:"runway"`
	Quality       []string               `json:"quality"`
}

type FlightPlanContract struct {
	FlightRules    string                          `json:"flightRules"`
	FlightType     string                          `json:"flightType"`
	RevisionNumber *int64                          `json:"revisionNo"`
	Status         string                          `json:"status"`
	Route          string                          `json:"route"`
	Altitude       *DistanceFlightPlanUnitContract `json:"altitude"`
	Airspeed       *SpeedFlightPlanUnitContract    `json:"airspeed"`
	LastUpdatedUTC time.Time                       `json:"lastUpdatedUtc"`
}

// SpeedFlightPlanUnitContract represents the speed flight plan unit contract
type SpeedFlightPlanUnitContract struct {
	Requested *Speed `json:"requested"`
	Assigned  *Speed `json:"assigned"`
}

// DistanceFlightPlanUnitContract represents the distance flight plan unit contract
type DistanceFlightPlanUnitContract struct {
	Requested *Measurement `json:"requested"`
	Assigned  *Measurement `json:"assigned"`
}

// StringFlightSearchItemContractSearchResultCollectionContract represents a string flight search item contract search result collection contract
type StringFlightSearchItemContractSearchResultCollectionContract struct {
	SearchBy string                     `json:"searchBy"`
	Count    int64                      `json:"count"`
	Items    []FlightSearchItemContract `json:"items"`
}

// FlightSearchItemContract represents a flight search item contract
type FlightSearchItemContract struct {
	Number string `json:"number"`
}

type AirportFidsContract struct {
	Departures []AirportFlightContract `json:"departures"`
	Arrivals   []AirportFlightContract `json:"arrivals"`
}

// Flight represents a simplified flight response.
type Flight struct {
	Number string `json:"number"`
	Status string `json:"status"`
}

// FlightScheduleParams represents optional query parameters shared across flight schedule endpoints.
type FlightScheduleParams struct {
	Direction      string `url:"direction,omitempty"`
	WithLeg        *bool  `url:"withLeg,omitempty"`
	WithCancelled  *bool  `url:"withCancelled,omitempty"`
	WithCodeshared *bool  `url:"withCodeshared,omitempty"`
	WithCargo      *bool  `url:"withCargo,omitempty"`
	WithPrivate    *bool  `url:"withPrivate,omitempty"`
	WithLocation   *bool  `url:"withLocation,omitempty"`
}

// FlightScheduleRelativeParams represents query parameters for a relative-time flight schedule request.
type FlightScheduleRelativeParams struct {
	OffsetMinutes   int    `url:"offsetMinutes"`
	DurationMinutes int    `url:"durationMinutes"`
	Direction       string `url:"direction,omitempty"`
	WithLeg         bool   `url:"withLeg,omitempty"`
	WithCancelled   bool   `url:"withCancelled,omitempty"`
	WithCodeshared  bool   `url:"withCodeshared,omitempty"`
	WithCargo       bool   `url:"withCargo,omitempty"`
	WithPrivate     bool   `url:"withPrivate,omitempty"`
	WithLocation    bool   `url:"withLocation,omitempty"`
}

// AirportFlightContract represents an airport flight contract
type AirportFlightContract struct {
	Movement        *FlightAirportMovementContract `json:"movement"`
	Departure       *FlightAirportMovementContract `json:"departure"`
	Arrival         *FlightAirportMovementContract `json:"arrival"`
	Number          string                         `json:"number"`
	CallSign        *string                        `json:"callSign"`
	Status          FlightStatus                   `json:"status"`
	CodeshareStatus CodeshareStatus                `json:"codeshareStatus"`
	IsCargo         bool                           `json:"isCargo"`
	Aircraft        *FlightAircraftContract        `json:"aircraft"`
	Airline         *FlightAirlineContract         `json:"airline"`
	Location        *FlightLocationContract        `json:"location"`
}

// FlightAircraftContract represents the aircraft contract for a flight
type FlightAircraftContract struct {
	Registration *string           `json:"reg"`
	ModeS        *string           `json:"modeS"`
	Model        *string           `json:"model"`
	Image        *ResourceContract `json:"image"`
}

// FlightLocationContract represents the location contract for a flight
type FlightLocationContract struct {
	PressureAltitude Measurement `json:"pressureAltitude"`
	Altitude         Measurement `json:"altitude"`
	Pressure         Pressure    `json:"pressure"`
	GroundSpeed      Speed       `json:"groundSpeed"`
	TrueTrack        Azimuth     `json:"trueTrack"`
	VsiFPM           *int64      `json:"vsiFpm"`
	ReportedAtUTC    time.Time   `json:"reportedAtUtc"`
	Latitude         float64     `json:"lat"`
	Longitude        float64     `json:"lon"`
}

// FlightSearchParams represents the search parameters for a flight
type FlightSearchParams struct {
	DateLocalRole     FlightDirection `url:"dateLocalRole,omitempty"`
	WithAircraftImage *bool           `url:"withAircraftImage,omitempty"`
	WithLocation      *bool           `url:"withLocation,omitempty"`
	WithFlightPlan    *bool           `url:"withFlightPlan,omitempty"`
}

// NewFlightSearchParams creates a new FlightSearchParams instance
func NewFlightSearchParams(dateLocalRole FlightDirection, withAircraftImage, withLocation, withFlightPlan bool) *FlightSearchParams {
	return &FlightSearchParams{
		DateLocalRole:     dateLocalRole,
		WithAircraftImage: &withAircraftImage,
		WithLocation:      &withLocation,
		WithFlightPlan:    &withFlightPlan,
	}
}

// FlightInformationDisplayParams represents the display parameters for a flight
type FlightInformationDisplayParams struct {
	Direction       *FlightDirection `url:"direction,omitempty"`
	WithLeg         *bool            `url:"withLeg,omitempty"`
	WithCancelled   *bool            `url:"withCancelled,omitempty"`
	WithCodeshared  *bool            `url:"withCodeshared,omitempty"`
	WithCargo       *bool            `url:"withCargo,omitempty"`
	WithPrivate     *bool            `url:"withPrivate,omitempty"`
	WithLocation    *bool            `url:"withLocation,omitempty"`
	OffsetMinutes   *int64           `url:"offsetMinutes,omitempty"`
	DurationMinutes *int64           `url:"durationMinutes,omitempty"`
}

// NewFlightInformationDisplayParams creates a new FlightInformationDisplayParams instance
func NewFlightInformationDisplayParams(direction *FlightDirection, withLeg, withCancelled, withCodeshared, withCargo, withPrivate, withLocation *bool, offsetMinutes, durationMinutes *int64) *FlightInformationDisplayParams {
	return &FlightInformationDisplayParams{
		Direction:       direction,
		WithLeg:         withLeg,
		WithCancelled:   withCancelled,
		WithCodeshared:  withCodeshared,
		WithCargo:       withCargo,
		WithPrivate:     withPrivate,
		WithLocation:    withLocation,
		OffsetMinutes:   offsetMinutes,
		DurationMinutes: durationMinutes,
	}
}

// FlightSearchByTermParams represents the search parameters for a flight by term
type FlightSearchByTermParams struct {
	Term  string `url:"q"`
	Limit *int64 `url:"limit,omitempty"`
}

// NewFlightSearchByTermParams creates a new FlightSearchByTermParams instance
func NewFlightSearchByTermParams(term string, limit *int64) *FlightSearchByTermParams {
	return &FlightSearchByTermParams{
		Term:  term,
		Limit: limit,
	}
}
