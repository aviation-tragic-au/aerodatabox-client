package client

import "github.com/deanstalker-stream/aerodatabox-client/model"

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

// EndpointHealthCheck represents the healthcheck endpoints.
type EndpointHealthCheck int

const (
	EndpointHealthCheckFeedStatusGeneral EndpointHealthCheck = iota
	EndpointHealthCheckFeedStatusData
	EndpointHealthCheckAirports
)

var endpointHealthCheck = map[EndpointHealthCheck]string{
	EndpointHealthCheckFeedStatusGeneral: "/health/services/feeds/{service}",
	EndpointHealthCheckFeedStatusData:    "/health/services/airports/{icao}/feeds",
	EndpointHealthCheckAirports:          "/health/services/feeds/{service}/airports",
}

// String returns the string representation of the EndpointHealthCheck value from the predefined mapping.
func (e EndpointHealthCheck) String() string {
	return endpointHealthCheck[e]
}

// Replace substitutes all {token} placeholders found in the endpoint URL template
// with the corresponding values from the provided map. An error is returned if
// any token present in the template has no matching key in the map.
func (e EndpointHealthCheck) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// GetDataFeedStatus retrieves the status of a data feed - FREE TIER
func (c *Client) GetDataFeedStatus(service string, withHttpCode bool) (*model.FeedServiceStatusContract, error) {
	url, err := EndpointHealthCheckFeedStatusGeneral.Replace(map[string]string{"service": service})
	if err != nil {
		return nil, err
	}

	if withHttpCode {
		url = url + "?withHttpCode=true"
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output *model.FeedServiceStatusContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetDataFeedStatusByICAO retrieves the status of a data feed for a specific airport - FREE TIER
func (c *Client) GetDataFeedStatusByICAO(icao string) (*model.AirportFeedServiceStatusContract, error) {
	url, err := EndpointHealthCheckFeedStatusData.Replace(map[string]string{"icao": icao})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output *model.AirportFeedServiceStatusContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetDataFeedSupportedAirports retrieves the supported airports for a data feed - FREE TIER
func (c *Client) GetDataFeedSupportedAirports(service string) (*model.StringCollectionContract, error) {
	url, err := EndpointHealthCheckAirports.Replace(map[string]string{"service": service})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output *model.StringCollectionContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}
