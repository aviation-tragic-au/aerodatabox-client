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

import "encoding/json"

// AircraftSearchByEnum represents the search by field for the aircraft search API
type AircraftSearchByEnum int

const (
	AircraftSearchByID AircraftSearchByEnum = iota
	AircraftSearchByRegistration
	AircraftSearchByICAO
)

var aircraftSearchBy = map[AircraftSearchByEnum]string{
	AircraftSearchByID:           "Id",
	AircraftSearchByRegistration: "Reg",
	AircraftSearchByICAO:         "Icao24",
}

var aircraftSearchByReverse = map[string]AircraftSearchByEnum{
	"Id":     AircraftSearchByID,
	"Reg":    AircraftSearchByRegistration,
	"Icao24": AircraftSearchByICAO,
}

// MarshalJSON implements the json.Marshaler interface.
func (s *AircraftSearchByEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(aircraftSearchBy[*s])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *AircraftSearchByEnum) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*s = aircraftSearchByReverse[str]
	return nil
}

// String implements the fmt.Stringer interface.
func (s *AircraftSearchByEnum) String() string {
	return aircraftSearchBy[*s]
}

// EngineType represents the engine type of an aircraft
type EngineType int

const (
	EngineTypeUnknown EngineType = iota
	EngineTypeJet
	EngineTypeTurboprop
	EngineTypePiston
)

var engineTypes = map[EngineType]string{
	EngineTypeUnknown:   "Unknown",
	EngineTypeJet:       "Jet",
	EngineTypeTurboprop: "Turboprop",
	EngineTypePiston:    "Piston",
}

var engineTypesReverse = map[string]EngineType{
	"Unknown":   EngineTypeUnknown,
	"Jet":       EngineTypeJet,
	"Turboprop": EngineTypeTurboprop,
	"Piston":    EngineTypePiston,
}

// MarshalJSON implements the json.Marshaler interface.
func (e *EngineType) MarshalJSON() ([]byte, error) {
	return json.Marshal(engineTypes[*e])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *EngineType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*e = engineTypesReverse[str]
	return nil
}

// String returns the string representation of the engine type
func (e *EngineType) String() string {
	return engineTypes[*e]
}
