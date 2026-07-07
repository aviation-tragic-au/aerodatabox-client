package model

// AirportParams represents the parameters for the airport API
type AirportParams struct {
	WithTime    *bool `url:"withTime,omitempty"`
	WithRunways *bool `url:"withRunways,omitempty"`
}

// NewAirportParams creates a new airport params instance
func NewAirportParams(withTime, withRunways *bool) *AirportParams {
	return &AirportParams{
		WithTime:    withTime,
		WithRunways: withRunways,
	}
}

// AirportLocationParams represents the parameters for the airport location API
type AirportLocationParams struct {
	Latitude           float64 `url:"lat"`
	Longitude          float64 `url:"lon"`
	Limit              int64   `url:"limit"`
	RadiusKM           int64   `url:"radiusKm"`
	WithFlightInfoOnly bool    `url:"withFlightInfoOnly,omitempty"`
}

// NewAirportLocationParams creates a new airport location params instance
func NewAirportLocationParams(latitude, longitude float64, limit, radiusKM int64, withFlightInfoOnly bool) *AirportLocationParams {
	return &AirportLocationParams{
		Latitude:           latitude,
		Longitude:          longitude,
		Limit:              limit,
		RadiusKM:           radiusKM,
		WithFlightInfoOnly: withFlightInfoOnly,
	}
}

// AirportIPParams represents the parameters for the airport IP API
type AirportIPParams struct {
	Q                  string `url:"q"`
	Limit              int64  `url:"limit"`
	RadiusKM           int64  `url:"radiusKm"`
	WithFlightInfoOnly bool   `url:"withFlightInfoOnly,omitempty"`
}

// NewAirportIPParams creates a new airport IP params instance
func NewAirportIPParams(q string, limit, radiusKM int64, withFlightInfoOnly bool) *AirportIPParams {
	return &AirportIPParams{
		Q:                  q,
		Limit:              limit,
		RadiusKM:           radiusKM,
		WithFlightInfoOnly: withFlightInfoOnly,
	}
}

// AirportSearchParams represents the parameters for the airport search API
type AirportSearchParams struct {
	Q                  string `url:"q"`
	Limit              *int64 `url:"limit,omitempty"`
	WithFlightInfoOnly bool   `url:"withFlightInfoOnly,omitempty"`
	WithSearchByCode   bool   `url:"withSearchByCode,omitempty"`
}

// NewAirportSearchParams creates a new airport search params instance
func NewAirportSearchParams(q string, limit *int64, withFlightInfoOnly, withSearchByCode bool) *AirportSearchParams {
	return &AirportSearchParams{
		Q:                  q,
		Limit:              limit,
		WithFlightInfoOnly: withFlightInfoOnly,
		WithSearchByCode:   withSearchByCode,
	}
}

// AirportContract represents an airport contract
type AirportContract struct {
	Icao             *string                   `json:"icao"`
	Iata             *string                   `json:"iata"`
	LocalCode        *string                   `json:"localCode"`
	ShortName        *string                   `json:"shortName"`
	FullName         string                    `json:"fullName"`
	MunicipalityName *string                   `json:"municipalityName"`
	Location         *Point                    `json:"location"`
	Elevation        *Measurement              `json:"elevation"`
	Country          *LocationCode             `json:"country"`
	Continent        *LocationCode             `json:"continent"`
	TimeZone         string                    `json:"timeZone"`
	URLs             *AirportURL               `json:"urls"`
	Runways          []*RunwayContract         `json:"runways"`
	CurrentTime      *AirportLocalTimeContract `json:"currentTime"`
}

// RunwayContract represents a runway
type RunwayContract struct {
	Name               string       `json:"name"`
	TrueHeading        float64      `json:"trueHdg"`
	Length             *Measurement `json:"length"`
	Width              *Measurement `json:"width"`
	IsClosed           bool         `json:"isClosed"`
	Location           *Point       `json:"location"`
	Surface            SurfaceType  `json:"surface"`
	DisplacedThreshold *Measurement `json:"displacedThreshold"`
	HasLighting        *bool        `json:"hasLighting"`
}

// LocationCode represents a location code
type LocationCode struct {
	Code string  `json:"code"`
	Name *string `json:"name"`
}

// AirportURL represents the URLs for an airport
type AirportURL struct {
	WebSite     *string `json:"webSite"`
	Wikipedia   *string `json:"wikipedia"`
	Twitter     *string `json:"twitter"`
	LiveATC     *string `json:"liveATC"`
	FlightRadar *string `json:"flightRadar"`
	GoogleMaps  *string `json:"googleMaps"`
}

// ListingAirportContract represents the information for an airport
type ListingAirportContract struct {
	ICAO             *string `json:"icao"`
	IATA             *string `json:"iata"`
	LocalCode        *string `json:"localCode"`
	Name             string  `json:"name"`
	ShortName        *string `json:"shortName"`
	MunicipalityName *string `json:"municipalityName"`
	Location         *Point  `json:"location"`
	CountryCode      string  `json:"countryCode"`
	TimeZone         *string `json:"timeZone"`
}

// Airport represents flight-context airport data including scheduling information.
type Airport struct {
	ICAO             *string `json:"icao"`
	IATA             *string `json:"iata"`
	Name             string  `json:"name"`
	ShortName        *string `json:"shortName"`
	MunicipalityName *string `json:"municipalityName"`
	Location         *Point  `json:"location"`
	CountryCode      string  `json:"countryCode"`
	TimeZone         *string `json:"timeZone"`
	ScheduledTime    *Time   `json:"scheduledTime"`
	ActualTime       *Time   `json:"actualTime"`
	Terminal         *string `json:"terminal"`
	Gate             *string `json:"gate"`
	Delay            *int64  `json:"delay"`
	Baggage          *string `json:"baggage"`
}

// AirportLocation represents a minimal airport identifier used in route and distance records.
type AirportLocation struct {
	ICAO *string `json:"icao"`
	IATA *string `json:"iata"`
	Name string  `json:"name"`
}

// AirportDelay represents delay statistics for an airport.
type AirportDelay struct {
	Airport                    ListingAirportContract `json:"airport"`
	From                       Time                   `json:"from"`
	To                         Time                   `json:"to"`
	DeparturesDelayInformation DelayInformation       `json:"departuresDelayInformation"`
	ArrivalsDelayInformation   DelayInformation       `json:"arrivalsDelayInformation"`
}

// StringListingAirportContractSearchResultCollectionContract represents the result of a string listing airport search
type StringListingAirportContractSearchResultCollectionContract struct {
	SearchBy string                   `json:"searchBy"`
	Count    int64                    `json:"count"`
	Items    []ListingAirportContract `json:"items"`
}

// GeoCoordinatesContractListingAirportContractSearchResultCollectionContract represents the result of a geo coordinates listing airport search
type GeoCoordinatesContractListingAirportContractSearchResultCollectionContract struct {
	SearchBy GeoCoordinatesContract   `json:"searchBy"`
	Count    int64                    `json:"count"`
	Items    []ListingAirportContract `json:"items"`
}

// GeoCoordinatesContract represents a geo coordinates
type GeoCoordinatesContract Point
