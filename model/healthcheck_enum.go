package model

import "encoding/json"

// FeedServiceStatus describes the status of the feed service
type FeedServiceStatus int

const (
	FeedServiceStatusDown FeedServiceStatus = iota
	FeedServiceStatusDegraded
	FeedServiceStatusOKPartial
	FeedServiceStatusOK
	FeedServiceStatusUnknown     = -2
	FeedServiceStatusUnavailable = -1
)

var statusStrings = map[FeedServiceStatus]string{
	FeedServiceStatusDown:        "Down",
	FeedServiceStatusOK:          "OK",
	FeedServiceStatusUnknown:     "Unknown",
	FeedServiceStatusUnavailable: "Unavailable",
	FeedServiceStatusDegraded:    "Degraded",
	FeedServiceStatusOKPartial:   "OKPartial",
}

var statusReverse = map[string]FeedServiceStatus{
	"Down":        FeedServiceStatusDown,
	"OK":          FeedServiceStatusOK,
	"Unknown":     FeedServiceStatusUnknown,
	"Unavailable": FeedServiceStatusUnavailable,
	"Degraded":    FeedServiceStatusDegraded,
	"OKPartial":   FeedServiceStatusOKPartial,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FeedServiceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(statusStrings[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FeedServiceStatus) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = statusReverse[str]
	return nil
}

// String returns the string representation of the feed service status
func (s *FeedServiceStatus) String() string {
	return statusStrings[*s]
}

// FeedServiceEnum represents the feed services
type FeedServiceEnum int

const (
	FeedServiceEnumFlightSchedules FeedServiceEnum = iota
	FeedServiceEnumFlightLiveUpdates
	FeedServiceEnumADSBUpdates
)

var feedServices = map[FeedServiceEnum]string{
	FeedServiceEnumFlightSchedules:   "FlightSchedules",
	FeedServiceEnumFlightLiveUpdates: "FlightLiveUpdates",
	FeedServiceEnumADSBUpdates:       "AdsbUpdates",
}

var feedServicesReverse = map[string]FeedServiceEnum{
	"FlightSchedules":   FeedServiceEnumFlightSchedules,
	"FlightLiveUpdates": FeedServiceEnumFlightLiveUpdates,
	"AdsbUpdates":       FeedServiceEnumADSBUpdates,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *FeedServiceEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(feedServices[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *FeedServiceEnum) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = feedServicesReverse[str]
	return nil
}

// String returns the string representation of the feed service
func (s *FeedServiceEnum) String() string {
	return feedServices[*s]
}
