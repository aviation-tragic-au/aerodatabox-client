package model

import "encoding/json"

// LicenseType represents the license type for image use
type LicenseType int

const (
	LicenseTypeAllRightsReserved LicenseType = iota
	LicenseTypeAttributionNonCommercialShareAlikeCC
	LicenseTypeAttributionNonCommercialCC
	LicenseTypeAttributionNonCommercialNoDerivativesCC
	LicenseTypeAttributionCC
	LicenseTypeAttributionNoDerivativesCC
	LicenseTypeAttributionShareAlikeCC
	LicenseTypeNoKnownCopyrightRestrictions
	LicenseTypeUnitedStatesGovernmentWork
	LicenseTypePublicDomainDedicationCC0
	LicenseTypePublicDomainMark
)

var licenseType = map[LicenseType]string{
	LicenseTypeAllRightsReserved:                       "AllRightsReserved",
	LicenseTypeAttributionNonCommercialShareAlikeCC:    "AttributionNoncommercialShareAlikeCC",
	LicenseTypeAttributionNonCommercialCC:              "AttributionNoncommercialCC",
	LicenseTypeAttributionNonCommercialNoDerivativesCC: "AttributionNoncommercialNoDerivativesCC",
	LicenseTypeAttributionCC:                           "AttributionCC",
	LicenseTypeAttributionNoDerivativesCC:              "AttributionNoDerivativesCC",
	LicenseTypeAttributionShareAlikeCC:                 "AttributionShareAlikeCC",
	LicenseTypeNoKnownCopyrightRestrictions:            "NoKnownCopyrightRestrictions",
	LicenseTypeUnitedStatesGovernmentWork:              "UnitedStatesGovernmentWork",
	LicenseTypePublicDomainDedicationCC0:               "PublicDomainDedicationCC0",
	LicenseTypePublicDomainMark:                        "PublicDomainMark",
}

var licenseTypeByString = map[string]LicenseType{
	"AllRightsReserved":                       LicenseTypeAllRightsReserved,
	"AttributionNoncommercialShareAlikeCC":    LicenseTypeAttributionNonCommercialShareAlikeCC,
	"AttributionNoncommercialCC":              LicenseTypeAttributionNonCommercialCC,
	"AttributionNoncommercialNoDerivativesCC": LicenseTypeAttributionNonCommercialNoDerivativesCC,
	"AttributionCC":                           LicenseTypeAttributionCC,
	"AttributionNoDerivativesCC":              LicenseTypeAttributionNoDerivativesCC,
	"AttributionShareAlikeCC":                 LicenseTypeAttributionShareAlikeCC,
	"NoKnownCopyrightRestrictions":            LicenseTypeNoKnownCopyrightRestrictions,
	"UnitedStatesGovernmentWork":              LicenseTypeUnitedStatesGovernmentWork,
	"PublicDomainDedicationCC0":               LicenseTypePublicDomainDedicationCC0,
	"PublicDomainMark":                        LicenseTypePublicDomainMark,
}

// MarshalJSON implements the json.Marshaler interface.
func (l *LicenseType) MarshalJSON() ([]byte, error) {
	return json.Marshal(licenseType[*l])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *LicenseType) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*l = licenseTypeByString[str]
	return nil
}

// String implements the fmt.Stringer interface.
func (l *LicenseType) String() string {
	return licenseType[*l]
}

// DayTime represents the day time of a flight
type DayTime int

const (
	DayTimeNight DayTime = iota
	DayTimeTwilightAstronomical
	DayTimeTwilightNautical
	DayTimeTwilightCivil
	DayTimeDay
	DayTimeBlueHour
	DayTimeGoldenHour
)

var dayTime = map[DayTime]string{
	DayTimeNight:                "Night",
	DayTimeTwilightAstronomical: "TwilightAstronomical",
	DayTimeTwilightNautical:     "TwilightNautical",
	DayTimeTwilightCivil:        "TwilightCivil",
	DayTimeDay:                  "Day",
	DayTimeBlueHour:             "BlueHour",
	DayTimeGoldenHour:           "GoldenHour",
}

var dayTimeByString = map[string]DayTime{
	"Night":                DayTimeNight,
	"TwilightAstronomical": DayTimeTwilightAstronomical,
	"TwilightNautical":     DayTimeTwilightNautical,
	"TwilightCivil":        DayTimeTwilightCivil,
	"Day":                  DayTimeDay,
	"BlueHour":             DayTimeBlueHour,
	"GoldenHour":           DayTimeGoldenHour,
}

// MarshalJSON implements the json.Marshaler interface.
func (d *DayTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dayTime[*d])
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *DayTime) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	*d = dayTimeByString[str]
	return nil
}

// String implements the fmt.Stringer interface.
func (d *DayTime) String() string {
	return dayTime[*d]
}
