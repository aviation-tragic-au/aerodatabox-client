package model

import "encoding/json"

// FlightSubscriptionSubjectType represents the subject type for the flight alert API
type FlightSubscriptionSubjectType int

const (
	FlightSubscriptionSubjectTypeFlightByNumber FlightSubscriptionSubjectType = iota
	FlightSubscriptionSubjectTypeFlightByAirportICAO
)

var flightSubscriptionSubjectType = map[FlightSubscriptionSubjectType]string{
	FlightSubscriptionSubjectTypeFlightByNumber:      "FlightByNumber",
	FlightSubscriptionSubjectTypeFlightByAirportICAO: "FlightByAirportIcao",
}

var flightSubscriptionSubjectTypeReverse = map[string]FlightSubscriptionSubjectType{
	"FlightByNumber":      FlightSubscriptionSubjectTypeFlightByNumber,
	"FlightByAirportIcao": FlightSubscriptionSubjectTypeFlightByAirportICAO,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightSubscriptionSubjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightSubscriptionSubjectType[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightSubscriptionSubjectType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightSubscriptionSubjectTypeReverse[str]
	return nil
}

// String returns the API name for this subject type.
func (s *FlightSubscriptionSubjectType) String() string {
	return flightSubscriptionSubjectType[*s]
}

// FlightSubscriptionBillingType represents the billing type for the flight alert API
type FlightSubscriptionBillingType int

const (
	FlightSubscriptionBillingTypeLifetimeBased FlightSubscriptionBillingType = iota
	FlightSubscriptionBillingTypeCreditBased
)

var flightSubscriptionBillingType = map[FlightSubscriptionBillingType]string{
	FlightSubscriptionBillingTypeLifetimeBased: "LifetimeBased",
	FlightSubscriptionBillingTypeCreditBased:   "CreditBased",
}

var flightSubscriptionBillingTypeReverse = map[string]FlightSubscriptionBillingType{
	"LifetimeBased": FlightSubscriptionBillingTypeLifetimeBased,
	"CreditBased":   FlightSubscriptionBillingTypeCreditBased,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FlightSubscriptionBillingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(flightSubscriptionBillingType[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FlightSubscriptionBillingType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = flightSubscriptionBillingTypeReverse[str]
	return nil
}

// String returns the API name for this billing type.
func (s *FlightSubscriptionBillingType) String() string {
	return flightSubscriptionBillingType[*s]
}
