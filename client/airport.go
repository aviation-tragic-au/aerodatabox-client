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
	"github.com/deanstalker-stream/aerodatabox-client/model"
	"github.com/google/go-querystring/query"
)

type EndpointAirport int

const (
	EndpointAirportByCode EndpointAirport = iota
	EndpointAirportRunways
	EndpointAirportSearchByLocation
	EndpointAirportSearchByIP
	EndpointAirportSearchByTerm
)

var endpointAirport = map[EndpointAirport]string{
	EndpointAirportByCode:           "/airports/{codeType}/{code}",
	EndpointAirportRunways:          "/airports/{codeType}/{code}/runways",
	EndpointAirportSearchByLocation: "/airports/search/location",
	EndpointAirportSearchByIP:       "/airports/search/ip",
	EndpointAirportSearchByTerm:     "/airports/search/term",
}

func (e EndpointAirport) String() string {
	return endpointAirport[e]
}

func (e EndpointAirport) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// GetAirportByCode retrieves an airport by its IATA or ICAO code. TIER 1
func (c *Client) GetAirportByCode(code, codeType string, withTime, withRunways *bool) (*model.Airport, error) {
	url, err := EndpointAirportByCode.Replace(map[string]string{
		"codeType": codeType,
		"code":     code,
	})
	if err != nil {
		return nil, err
	}

	params := model.NewAirportParams(withTime, withRunways)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.Airport
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAirportRunways retrieves the runways of an airport - TIER 1
func (c *Client) GetAirportRunways(code, codeType string) ([]*model.RunwayContract, error) {
	url, err := EndpointAirportRunways.Replace(map[string]string{
		"codeType": codeType,
		"code":     code,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output []*model.RunwayContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAirportByLocation retrieves airports within a given radius of a location - TIER 2
func (c *Client) GetAirportByLocation(latitude, longitude float64, limit, radiusKm int64, withFlightInfoOnly bool) (*model.StringListingAirportContractSearchResultCollectionContract, error) {
	url := EndpointAirportSearchByLocation.String()

	params := model.NewAirportLocationParams(latitude, longitude, limit, radiusKm, withFlightInfoOnly)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.StringListingAirportContractSearchResultCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAirportByIP retrieves airports within a given radius of an IP address - TIER 2
func (c *Client) GetAirportByIP(q string, limit, radiusKM int64, withFlightInfoOnly bool) (*model.StringListingAirportContractSearchResultCollectionContract, error) {
	url := EndpointAirportSearchByIP.String()

	params := model.NewAirportIPParams(q, limit, radiusKM, withFlightInfoOnly)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.StringListingAirportContractSearchResultCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAirportByTerm retrieves airports by a search term - TIER 2
func (c *Client) GetAirportByTerm(q string, limit *int64, withFlightInfoOnly, withSearchByCode bool) (*model.StringListingAirportContractSearchResultCollectionContract, error) {
	url := EndpointAirportSearchByTerm.String()

	params := model.NewAirportSearchParams(q, limit, withFlightInfoOnly, withSearchByCode)
	v, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.StringListingAirportContractSearchResultCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}
