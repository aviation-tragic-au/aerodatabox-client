package client

import (
	"fmt"
	"time"

	"github.com/deanstalker-stream/aerodatabox-client/model"
	"github.com/google/go-querystring/query"
)

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

// EndpointFlight represents the flight endpoints.
type EndpointFlight int

const (
	EndpointFlightStatus EndpointFlight = iota
	EndpointFlightStatusByDate
	EndpointFlightScheduleHistory
	EndpointFlightDepartureDates
	EndpointFlightDepartureDateRange
	EndpointFlightFIDSLocalTime
	EndpointFlightFIDSRelativeTime
	EndpointFlightSearchByTerm
)

var endpointFlight = map[EndpointFlight]string{
	EndpointFlightStatus:             "/flights/{searchBy}/{searchParameter}",
	EndpointFlightStatusByDate:       "/flights/{searchBy}/{searchParameter}/{dateLocal}",
	EndpointFlightScheduleHistory:    "/flights/{searchBy}/{searchParameter}/{dateFromLocal}/{dateToLocal}",
	EndpointFlightDepartureDates:     "/flights/{searchBy}/{searchParameter}/dates",
	EndpointFlightDepartureDateRange: "/flights/{searchBy}/{searchParameter}/dates/{fromLocal}/{toLocal}",
	EndpointFlightFIDSLocalTime:      "/flights/airports/{codeType}/{code}/{fromLocal}/{toLocal}",
	EndpointFlightFIDSRelativeTime:   "/flights/airports/{codeType}/{code}",
	EndpointFlightSearchByTerm:       "/flights/search/term",
}

// String returns the string representation of the EndpointFlight.
func (e EndpointFlight) String() string {
	return endpointFlight[e]
}

// Replace replaces tokens in the endpoint template with values provided in the tokens map and returns the processed string.
// Returns an error if any required token is missing from the tokens map.
func (e EndpointFlight) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

func validateFlightSearchParam(s string) error {
	if len(s) > 10 {
		return fmt.Errorf("searchParameter %q exceeds maximum length of 10", s)
	}
	return nil
}

// GetFlightStatus retrieves the status of a flight - TIER 2
func (c *Client) GetFlightStatus(searchBy, searchParameter string, dateLocalRole model.FlightDirection, withAircraftImage, withLocation, withFlightPlan bool) ([]*model.FlightContract, error) {
	if err := validateFlightSearchParam(searchParameter); err != nil {
		return nil, err
	}

	url, err := EndpointFlightStatus.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
	})
	if err != nil {
		return nil, err
	}

	params := model.NewFlightSearchParams(dateLocalRole, withAircraftImage, withLocation, withFlightPlan)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(appendQuery(url, v.Encode()))
	if err != nil {
		return nil, err
	}

	var output []*model.FlightContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightStatusByDate retrieves the status of a flight for a specific date - TIER 2
func (c *Client) GetFlightStatusByDate(searchBy, searchParameter string, dateLocal time.Time, dateLocalRole model.FlightDirection, withAircraftImage, withLocation, withFlightPlan bool) ([]*model.FlightContract, error) {
	if err := validateFlightSearchParam(searchParameter); err != nil {
		return nil, err
	}

	url, err := EndpointFlightStatusByDate.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
		"dateLocal":       dateLocal.Format("2006-01-02T15:04"),
	})
	if err != nil {
		return nil, err
	}

	params := model.NewFlightSearchParams(dateLocalRole, withAircraftImage, withLocation, withFlightPlan)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(appendQuery(url, v.Encode()))
	if err != nil {
		return nil, err
	}

	var output []*model.FlightContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightScheduleHistory retrieves the schedule history of a flight - TIER 2
func (c *Client) GetFlightScheduleHistory(searchBy, searchParameter string, dateFromLocal, dateToLocal time.Time, dateLocalRole *model.FlightDirection) ([]*model.FlightContract, error) {
	if err := validateFlightSearchParam(searchParameter); err != nil {
		return nil, err
	}

	url, err := EndpointFlightScheduleHistory.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
		"dateFromLocal":   dateFromLocal.Format("2006-01-02T15:04"),
		"dateToLocal":     dateToLocal.Format("2006-01-02T15:04"),
	})
	if err != nil {
		return nil, err
	}

	if dateLocalRole != nil {
		params := model.NewFlightSearchParams(*dateLocalRole, false, false, false)
		v, err := query.Values(params)
		if err != nil {
			return nil, err
		}
		url = appendQuery(url, v.Encode())
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output []*model.FlightContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightDepartureDates retrieves the departure dates of a flight - TIER 2
func (c *Client) GetFlightDepartureDates(searchBy, searchParameter string) ([]string, error) {
	if err := validateFlightSearchParam(searchParameter); err != nil {
		return nil, err
	}

	url, err := EndpointFlightDepartureDates.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output []string
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightDepartureDateRange retrieves the departure dates of a flight within a specific date range - TIER 2
func (c *Client) GetFlightDepartureDateRange(searchBy, searchParameter string, fromLocal, toLocal time.Time) ([]string, error) {
	if err := validateFlightSearchParam(searchParameter); err != nil {
		return nil, err
	}

	url, err := EndpointFlightDepartureDateRange.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
		"fromLocal":       fromLocal.Format("2006-01-02T15:04"),
		"toLocal":         toLocal.Format("2006-01-02T15:04"),
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output []string
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightInformationDepartureByLocalTime retrieves the departure information of a flight based on a local time - TIER 2
func (c *Client) GetFlightInformationDepartureByLocalTime(codeType, code string, fromLocal, toLocal time.Time, direction *model.FlightDirection, withLeg, withCancelled, withCodeshared, withCargo, withPrivate, withLocation *bool) (*model.AirportFidsContract, error) {
	url, err := EndpointFlightFIDSLocalTime.Replace(map[string]string{
		"codeType":  codeType,
		"code":      code,
		"fromLocal": fromLocal.Format("2006-01-02T15:04"),
		"toLocal":   toLocal.Format("2006-01-02T15:04"),
	})
	if err != nil {
		return nil, err
	}

	params := model.NewFlightInformationDisplayParams(direction, withLeg, withCancelled, withCodeshared, withCargo, withPrivate, withLocation, nil, nil)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(appendQuery(url, v.Encode()))
	if err != nil {
		return nil, err
	}

	var output *model.AirportFidsContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightInformationDepartureByRelativeTime retrieves the departure information of a flight based on a relative time offset - TIER 2
func (c *Client) GetFlightInformationDepartureByRelativeTime(codeType, code string, direction *model.FlightDirection, withLeg, withCancelled, withCodeshared, withCargo, withPrivate, withLocation *bool, offsetMinutes, durationMinutes *int64) (*model.AirportFidsContract, error) {
	url, err := EndpointFlightFIDSRelativeTime.Replace(map[string]string{
		"codeType": codeType,
		"code":     code,
	})
	if err != nil {
		return nil, err
	}

	params := model.NewFlightInformationDisplayParams(direction, withLeg, withCancelled, withCodeshared, withCargo, withPrivate, withLocation, offsetMinutes, durationMinutes)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(appendQuery(url, v.Encode()))
	if err != nil {
		return nil, err
	}

	var output *model.AirportFidsContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// SearchFlightByTerm searches for flights by a given search term - TIER 2
func (c *Client) SearchFlightByTerm(term string, limit int64) (*model.StringFlightSearchItemContractSearchResultCollectionContract, error) {
	url := EndpointFlightSearchByTerm.String()

	params := model.NewFlightSearchByTermParams(term, &limit)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(appendQuery(url, v.Encode()))
	if err != nil {
		return nil, err
	}

	var output *model.StringFlightSearchItemContractSearchResultCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}
