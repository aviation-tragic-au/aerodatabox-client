package model

import "encoding/json"

// AirportCodesByEnum represents the airport codes by field for the airport API
type AirportCodesByEnum int

const (
	AirportCodesByICAO AirportCodesByEnum = iota
	AirportCodesByIATA
)

var airportCodesBy = map[AirportCodesByEnum]string{
	AirportCodesByICAO: "Icao",
	AirportCodesByIATA: "Iata",
}

var airportCodesByReverse = map[string]AirportCodesByEnum{
	"Icao": AirportCodesByICAO,
	"Iata": AirportCodesByIATA,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *AirportCodesByEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(airportCodesBy[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *AirportCodesByEnum) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = airportCodesByReverse[str]
	return nil
}

// String implements the fmt.Stringer interface.
func (s *AirportCodesByEnum) String() string {
	return airportCodesBy[*s]
}

// SurfaceType represents the surface type of airport
type SurfaceType int

const (
	SurfaceTypeUnknown SurfaceType = iota
	SurfaceTypeAsphalt
	SurfaceTypeConcrete
	SurfaceTypeGrass
	SurfaceTypeDirt
	SurfaceTypeGravel
	SurfaceTypeDryLakebed
	SurfaceTypeWater
	SurfaceTypeSnow
)

var surfaceTypes = map[SurfaceType]string{
	SurfaceTypeUnknown:    "Unknown",
	SurfaceTypeAsphalt:    "Asphalt",
	SurfaceTypeConcrete:   "Concrete",
	SurfaceTypeGrass:      "Grass",
	SurfaceTypeDirt:       "Dirt",
	SurfaceTypeGravel:     "Gravel",
	SurfaceTypeDryLakebed: "DryLakebed",
	SurfaceTypeWater:      "Water",
	SurfaceTypeSnow:       "Snow",
}

var surfaceTypesByReverse = map[string]SurfaceType{
	"Unknown":    SurfaceTypeUnknown,
	"Asphalt":    SurfaceTypeAsphalt,
	"Concrete":   SurfaceTypeConcrete,
	"Grass":      SurfaceTypeGrass,
	"Dirt":       SurfaceTypeDirt,
	"Gravel":     SurfaceTypeGravel,
	"DryLakebed": SurfaceTypeDryLakebed,
	"Water":      SurfaceTypeWater,
	"Snow":       SurfaceTypeSnow,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *SurfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(surfaceTypes[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *SurfaceType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = surfaceTypesByReverse[str]
	return nil
}

// String returns the string representation of the surface type
func (s *SurfaceType) String() string {
	return surfaceTypes[*s]
}
