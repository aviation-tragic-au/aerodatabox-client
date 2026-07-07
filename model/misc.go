package model

// This file contains models related to the Miscellaneous API endpoints for Aerodatabox

import (
	"fmt"
	"net/url"
)

// Pressure represents a pressure
type Pressure struct {
	Hectopascal   float64 `json:"hPa"`
	Millimeter    float64 `json:"mmHg"`
	InchOfMercury float64 `json:"inHg"`
}

// Time represents a time
type Time struct {
	UTC   string `json:"utc"`
	Local string `json:"local"`
}

// Timezone represents a timezone
type Timezone struct {
	TimeZone string `json:"timeZone"`
}

// Angle represents an angle
type Angle struct {
	Degree float64 `json:"deg"`
	Radian float64 `json:"rad"`
}

// Azimuth represents an azimuth
type Azimuth struct {
	Degree  float64 `json:"deg"`
	Radians float64 `json:"rad"`
}

// Point represents a point
type Point struct {
	Latitude  *float64 `json:"lat,omitempty"`
	Longitude *float64 `json:"lon,omitempty"`
}

// AirportLocalTimeContract response for /airports/<codeType>/<code>/time/local
type AirportLocalTimeContract struct {
	Time       Time   `json:"time"`
	TimeZoneID string `json:"timeZoneId"`
}

// Speed represents a speed
type Speed struct {
	Knots float64 `json:"kt"`
	Kmh   float64 `json:"kmPerHour"`
	Mph   float64 `json:"miPerHour"`
}

// DayCycle represents a day cycle
type DayCycle struct {
	Location         Point    `json:"location"`
	SunElevation     Angle    `json:"sunElevation"`
	SunAzimuth       Angle    `json:"sunAzimuth"`
	DayTime          []string `json:"dayTime"`
	DawnAstronomical Time     `json:"dawnAstronomical"`
	DawnNautical     Time     `json:"dawnNautical"`
	DawnCivil        Time     `json:"dawnCivil"`
	Sunrise          Time     `json:"sunrise"`
	NoonTrue         Time     `json:"noonTrue"`
	Sunset           Time     `json:"sunset"`
	DuskCivil        Time     `json:"duskCivil"`
	DuskNautical     Time     `json:"duskNautical"`
	DuskAstronomical Time     `json:"duskAstronomical"`
}

type FlightTimeParams struct {
	flightTimeModel *string
	aircraftName    *string
}

// NewFlightTimeParams creates a new FlightTimeParams instance
func NewFlightTimeParams(flightTimeModel, aircraftName *string) (*FlightTimeParams, error) {
	if flightTimeModel != nil {
		if *flightTimeModel != "Standard" && *flightTimeModel != "ML01" {
			return nil, fmt.Errorf("flightTimeModel must be either Standard or ML01")
		}
	}
	return &FlightTimeParams{
		flightTimeModel: flightTimeModel,
		aircraftName:    aircraftName,
	}, nil
}

// EncodeValues encodes the parameters into a url.Values object
func (p *FlightTimeParams) EncodeValues(_ string, v *url.Values) error {
	if p.flightTimeModel != nil {
		v.Set("flightTimeModel", *p.flightTimeModel)
	}
	if p.aircraftName != nil {
		v.Set("aircraftName", *p.aircraftName)
	}
	return nil
}

// AirportDistanceTimeContract response for /airports/<codeType>/<codeFrom>/distance-time/<codeTo>
type AirportDistanceTimeContract struct {
	From                  AirportLocation     `json:"from"`
	To                    AirportLocation     `json:"to"`
	GreatCircleDistance   GreatCircleDistance `json:"greatCircleDistance"`
	ApproximateFlightTime string              `json:"approxFlightTime"`
}

// GreatCircleDistance represents a great circle distance
type GreatCircleDistance struct {
	Meter        float64 `json:"meter"`
	Kilometre    float64 `json:"km"`
	Mile         float64 `json:"mile"`
	NauticalMile float64 `json:"nm"`
	Feet         float64 `json:"feet"`
}

// Measurement represents an altitude
type Measurement struct {
	Meter float64 `json:"meter"`
	KM    float64 `json:"km"`
	Mile  float64 `json:"mile"`
	NM    float64 `json:"nm"`
	Feet  float64 `json:"feet"`
}

// Track represents a track
type Track struct {
	Degrees float64 `json:"deg"`
	Radians float64 `json:"rad"`
}

// Altitude represents an altitude or speed measurement value with a unit.
type Altitude struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// APIError represents an error response from the AeroDataBox API
type APIError struct {
	StatusCode int     `json:"-"`
	Message    string  `json:"message"`
	Details    *string `json:"details"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != nil && *e.Details != "" {
		return fmt.Sprintf("aerodatabox API error %d: %s (%s)", e.StatusCode, e.Message, *e.Details)
	}
	return fmt.Sprintf("aerodatabox API error %d: %s", e.StatusCode, e.Message)
}

// SolarStateContract response for /solar-state
type SolarStateContract struct {
	Location         Point     `json:"location"`
	SunElevation     Angle     `json:"sunElevation"`
	SunAzimuth       Azimuth   `json:"sunAzimuth"`
	DayTime          []DayTime `json:"dayTime"`
	DawnAstronomical *Time     `json:"dawnAstronomical"`
	DawnCivil        *Time     `json:"dawnCivil"`
	Sunrise          *Time     `json:"sunrise"`
	NoonTrue         *Time     `json:"noonTrue"`
	Sunset           *Time     `json:"sunset"`
	DuskCivil        *Time     `json:"duskCivil"`
	DuskNautical     *Time     `json:"duskNautical"`
	DuskAstronomical *Time     `json:"duskAstronomical"`
}
