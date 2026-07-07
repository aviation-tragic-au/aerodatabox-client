package model

import "encoding/json"

// StatisticClass represents the statistic class
type StatisticClass int

const (
	StatisticClassFlight StatisticClass = iota
	StatisticClassFlightAndHour
)

var statisticClass = map[StatisticClass]string{
	StatisticClassFlight:        "Flight",
	StatisticClassFlightAndHour: "FlightAndHour",
}

var statisticClassReverse = map[string]StatisticClass{
	"Flight":        StatisticClassFlight,
	"FlightAndHour": StatisticClassFlightAndHour,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *StatisticClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(statisticClass[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *StatisticClass) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = statisticClassReverse[str]
	return nil
}

// String returns the string representation of the statistic class
func (s *StatisticClass) String() string {
	return statisticClass[*s]
}
