package model

import (
	"encoding/json"
	"net/url"
)

// CodeshareStatus represents the codeshare status for a flight
type CodeshareStatus int

const (
	CodeshareStatusUnknown CodeshareStatus = iota
	CodeshareStatusIsOperator
	CodeshareStatusIsCodeshared
)

var codeshareStatus = map[CodeshareStatus]string{
	CodeshareStatusUnknown:      "Unknown",
	CodeshareStatusIsOperator:   "IsOperator",
	CodeshareStatusIsCodeshared: "IsCodeshared",
}

var codeshareStatusReverse = map[string]CodeshareStatus{
	"Unknown":      CodeshareStatusUnknown,
	"IsOperator":   CodeshareStatusIsOperator,
	"IsCodeshared": CodeshareStatusIsCodeshared,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *CodeshareStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(codeshareStatus[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *CodeshareStatus) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = codeshareStatusReverse[str]
	return nil
}

// String returns the string representation of the codeshare status
func (s *CodeshareStatus) String() string {
	return codeshareStatus[*s]
}

// FlightDirection represents the direction of a flight
type FlightDirection int

const (
	FlightDirectionArrival   FlightDirection = 1
	FlightDirectionDeparture FlightDirection = 2
	FlightDirectionBoth      FlightDirection = 3
)

var flightDirection = map[FlightDirection]string{
	FlightDirectionArrival:   "Arrival",
	FlightDirectionDeparture: "Departure",
	FlightDirectionBoth:      "Both",
}

var flightDirectionReverse = map[string]FlightDirection{
	"Arrival":   FlightDirectionArrival,
	"Departure": FlightDirectionDeparture,
	"Both":      FlightDirectionBoth,
}

// EncodeValues implements go-querystring's Encoder interface so FlightDirection is encoded as its string label.
func (s FlightDirection) EncodeValues(key string, v *url.Values) error {
	if str := flightDirection[s]; str != "" {
		v.Set(key, str)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightDirection) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightDirection[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightDirection) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightDirectionReverse[str]
	return nil
}

// String returns the string representation of the flight direction
func (s *FlightDirection) String() string {
	return flightDirection[*s]
}

// FlightPlanStatus represents the status of a flight plan
type FlightPlanStatus int

const (
	FlightPlanStatusProposed FlightPlanStatus = iota
	FlightPlanStatusActive
	FlightPlanStatusDropped
	FlightPlanStatusCancelled
	FlightPlanStatusCompleted
)

var flightPlanStatus = map[FlightPlanStatus]string{
	FlightPlanStatusProposed:  "Proposed",
	FlightPlanStatusActive:    "Active",
	FlightPlanStatusDropped:   "Dropped",
	FlightPlanStatusCancelled: "Cancelled",
	FlightPlanStatusCompleted: "Completed",
}

var flightPlanStatusReverse = map[string]FlightPlanStatus{
	"Proposed":  FlightPlanStatusProposed,
	"Active":    FlightPlanStatusActive,
	"Dropped":   FlightPlanStatusDropped,
	"Cancelled": FlightPlanStatusCancelled,
	"Completed": FlightPlanStatusCompleted,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightPlanStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightPlanStatus[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightPlanStatus) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightPlanStatusReverse[str]
	return nil
}

// String returns the string representation of the flight plan status
func (s *FlightPlanStatus) String() string {
	return flightPlanStatus[*s]
}

type FlightRules int

const (
	FlightRulesVFR FlightRules = iota
	FlightRulesIFR
)

var flightRules = map[FlightRules]string{
	FlightRulesVFR: "VFR",
	FlightRulesIFR: "IFR",
}

var flightRulesReverse = map[string]FlightRules{
	"VFR": FlightRulesVFR,
	"IFR": FlightRulesIFR,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightRules) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightRules[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightRules) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightRulesReverse[str]
	return nil
}

// String returns the string representation of the flight rules
func (s *FlightRules) String() string {
	return flightRules[*s]
}

// FlightSearchByEnum represents the flight search by enum
type FlightSearchByEnum int

const (
	FlightSearchByFlightNumber FlightSearchByEnum = iota
	FlightSearchByRegistration
	FlightSearchByCallSign
	FlightSearchByIcao24
)

var flightSearchBy = map[FlightSearchByEnum]string{
	FlightSearchByFlightNumber: "Number",
	FlightSearchByRegistration: "Reg",
	FlightSearchByCallSign:     "CallSign",
	FlightSearchByIcao24:       "Icao24",
}

var flightSearchByReverse = map[string]FlightSearchByEnum{
	"Number":   FlightSearchByFlightNumber,
	"Reg":      FlightSearchByRegistration,
	"CallSign": FlightSearchByCallSign,
	"Icao24":   FlightSearchByIcao24,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightSearchByEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightSearchBy[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightSearchByEnum) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightSearchByReverse[str]
	return nil
}

// String returns the string representation of the flight search by enum
func (s *FlightSearchByEnum) String() string {
	return flightSearchBy[*s]
}

// FlightStatus represents the status of a flight
type FlightStatus int

const (
	FlightStatusUnknown FlightStatus = iota
	FlightStatusExpected
	FlightStatusEnRoute
	FlightStatusCheckin
	FlightStatusBoarding
	FlightStatusGateClosed
	FlightStatusDeparted
	FlightStatusDelayed
	FlightStatusApproaching
	FlightStatusArrived
	FlightStatusCanceled
	FlightStatusDiverted
	FlightStatusCanceledUncertain
)

var flightStatus = map[FlightStatus]string{
	FlightStatusUnknown:           "Unknown",
	FlightStatusExpected:          "Expected",
	FlightStatusEnRoute:           "EnRoute",
	FlightStatusCheckin:           "CheckIn",
	FlightStatusBoarding:          "Boarding",
	FlightStatusGateClosed:        "GateClosed",
	FlightStatusDeparted:          "Departed",
	FlightStatusDelayed:           "Delayed",
	FlightStatusApproaching:       "Approaching",
	FlightStatusArrived:           "Arrived",
	FlightStatusCanceled:          "Canceled",
	FlightStatusDiverted:          "Diverted",
	FlightStatusCanceledUncertain: "CanceledUncertain",
}

var flightStatusReverse = map[string]FlightStatus{
	"Unknown":           FlightStatusUnknown,
	"Expected":          FlightStatusExpected,
	"EnRoute":           FlightStatusEnRoute,
	"CheckIn":           FlightStatusCheckin,
	"Boarding":          FlightStatusBoarding,
	"GateClosed":        FlightStatusGateClosed,
	"Departed":          FlightStatusDeparted,
	"Delayed":           FlightStatusDelayed,
	"Approaching":       FlightStatusApproaching,
	"Arrived":           FlightStatusArrived,
	"Canceled":          FlightStatusCanceled,
	"Diverted":          FlightStatusDiverted,
	"CanceledUncertain": FlightStatusCanceledUncertain,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightStatus[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightStatus) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightStatusReverse[str]
	return nil
}

// String returns the string representation of the flight status
func (s *FlightStatus) String() string {
	return flightStatus[*s]
}

// FlightType represents the type of a flight
type FlightType int

const (
	FlightTypeOther FlightType = iota
	FlightTypeGeneral
	FlightTypeScheduled
	FlightTypeNonScheduled
	FlightTypeMilitary
)

var flightType = map[FlightType]string{
	FlightTypeOther:        "Other",
	FlightTypeGeneral:      "General",
	FlightTypeScheduled:    "Scheduled",
	FlightTypeNonScheduled: "NonScheduled",
	FlightTypeMilitary:     "Military",
}

var flightTypeReverse = map[string]FlightType{
	"Other":        FlightTypeOther,
	"General":      FlightTypeGeneral,
	"Scheduled":    FlightTypeScheduled,
	"NonScheduled": FlightTypeNonScheduled,
	"Military":     FlightTypeMilitary,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightType) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightType[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightTypeReverse[str]
	return nil
}

// String returns the string representation of the flight type
func (s *FlightType) String() string {
	return flightType[*s]
}

// FlightTimeEnum represents the flight time of a model
type FlightTimeEnum int

const (
	FlightTimeEnumStandard FlightTimeEnum = iota
	FlightTimeEnumML01
)

var modelFlightTime = map[FlightTimeEnum]string{
	FlightTimeEnumStandard: "Standard",
	FlightTimeEnumML01:     "ML01",
}

var modelFlightTimeReverse = map[string]FlightTimeEnum{
	"Standard": FlightTimeEnumStandard,
	"ML01":     FlightTimeEnumML01,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightTimeEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(modelFlightTime[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightTimeEnum) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = modelFlightTimeReverse[str]
	return nil
}

// String returns the string representation of the flight time
func (s *FlightTimeEnum) String() string {
	return modelFlightTime[*s]
}
