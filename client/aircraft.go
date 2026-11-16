package client

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

import (
	"github.com/aviation-tragic-au/aerodatabox-client/model"
	"github.com/google/go-querystring/query"
)

type EndpointAircraft int

const (
	EndpointAircraftSearchSingle EndpointAircraft = iota
	EndpointAircraftSearchAll
	EndpointAircraftRegoHistory
	EndpointAircraftSearchByTerm
	EndpointAirlineFleetByCode
	EndpointAircraftImageByTailNumber
)

var endpointAircraft = map[EndpointAircraft]string{
	EndpointAircraftSearchSingle:      "/aircrafts/{searchBy}/{searchParameter}",               //nolint:misspell
	EndpointAircraftSearchAll:         "/aircrafts/{searchBy}/{searchParameter}/all",           //nolint:misspell
	EndpointAircraftRegoHistory:       "/aircrafts/{searchBy}/{searchParameter}/registrations", //nolint:misspell
	EndpointAircraftSearchByTerm:      "/aircrafts/search/term",                                //nolint:misspell
	EndpointAirlineFleetByCode:        "/airlines/{code}/aircrafts",                            //nolint:misspell
	EndpointAircraftImageByTailNumber: "/aircrafts/reg/{reg}/image/beta",                       //nolint:misspell
}

func (e EndpointAircraft) String() string {
	return endpointAircraft[e]
}

// Replace substitutes all {token} placeholders found in the endpoint URL template
// with the corresponding values from the provided map. An error is returned if
// any token present in the template has no matching key in the map.
func (e EndpointAircraft) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// GetSingleAircraft retrieves a single aircraft by its registration or ICAO code - TIER 1
func (c *Client) GetSingleAircraft(searchBy, searchParameter string, withRegistrations, withImage *bool) (*model.AircraftContract, error) {
	url, err := EndpointAircraftSearchSingle.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
	})
	if err != nil {
		return nil, err
	}

	params := model.NewAircraftParams(withRegistrations, withImage)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.AircraftContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAircraft retrieves a list of aircraft by their registration or ICAO code.
func (c *Client) GetAircraft(searchBy, searchParameter string, withRegistrations, withImage *bool) ([]*model.AircraftContract, error) {
	url, err := EndpointAircraftSearchAll.Replace(map[string]string{
		"searchBy":        searchBy,
		"searchParameter": searchParameter,
	})
	if err != nil {
		return nil, err
	}
	params := model.NewAircraftParams(withRegistrations, withImage)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output []*model.AircraftContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAircraftHistory retrieves the registration history of an aircraft.
func (c *Client) GetAircraftHistory(searchBy, searchParameter string) (*model.AircraftHistory, error) {
	url, err := EndpointAircraftRegoHistory.Replace(map[string]string{
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

	var output *model.AircraftHistory
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAircraftByTerm retrieves aircraft by a search term.
func (c *Client) GetAircraftByTerm(term string, limit *int64) (*model.StringAircraftContractSearchResultCollectionContract, error) {
	url := EndpointAircraftSearchByTerm.String()

	params := model.NewAircraftSearchParams(term, limit)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.StringAircraftContractSearchResultCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAirlineFleetByCode retrieves the fleet of an airline by its IATA code.
func (c *Client) GetAirlineFleetByCode(code string, pageSize, pageOffset *int64, withRegistrations *bool) (*model.AircraftContractPagedCollectionContract, error) {
	url, err := EndpointAirlineFleetByCode.Replace(map[string]string{
		"code": code,
	})
	if err != nil {
		return nil, err
	}

	params := model.NewAirlineFleetByCodeParams(pageSize, pageOffset, withRegistrations)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.AircraftContractPagedCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}
