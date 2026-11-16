package client

import (
	"context"
	"fmt"
	neturl "net/url"
	"strings"
	"time"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
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

// EndpointMisc represents the misc endpoints.
type EndpointMisc int

const (
	EndpointMiscAirportCurrentTime EndpointMisc = iota
	EndpointMiscAirportSolarTime
	EndpointMiscAirportSolarTimeByDate
	EndpointMiscFlightTime
)

var endpointMisc = map[EndpointMisc]string{
	EndpointMiscAirportCurrentTime:     "/airports/{codeType}/{code}/time/local",
	EndpointMiscAirportSolarTime:       "/airports/{codeType}/{code}/time/solar",
	EndpointMiscAirportSolarTimeByDate: "/airports/{codeType}/{code}/time/solar/{dateLocal}",
	EndpointMiscFlightTime:             "/airports/{codeType}/{codeFrom}/distance-time/{codeTo}",
}

// String returns the endpoint URL.
func (e EndpointMisc) String() string {
	return endpointMisc[e]
}

// Replace returns the endpoint with tokens replaced.
func (e EndpointMisc) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// GetAirportCurrentTime retrieves the current time of an airport - TIER 1
func (c *Client) GetAirportCurrentTime(ctx context.Context, code, codeType string) (*model.AirportLocalTimeContract, error) {
	if err := c.validateCode(codeType, code); err != nil {
		return nil, err
	}
	url, err := EndpointMiscAirportCurrentTime.Replace(map[string]string{
		"codeType": codeType,
		"code":     code,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var output *model.AirportLocalTimeContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetAirportSolarTime retrieves the solar time of an airport - TIER 1
func (c *Client) GetAirportSolarTime(ctx context.Context, code, codeType string) (*model.SolarStateContract, error) {
	if err := c.validateCode(codeType, code); err != nil {
		return nil, err
	}
	url, err := EndpointMiscAirportSolarTime.Replace(map[string]string{
		"codeType": codeType,
		"code":     code,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var output *model.SolarStateContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetAirportSolarTimeByDate retrieves the solar time of an airport at a specific date - TIER 1
func (c *Client) GetAirportSolarTimeByDate(ctx context.Context, code, codeType string, dateLocal time.Time) (*model.SolarStateContract, error) {
	if err := c.validateCode(codeType, code); err != nil {
		return nil, err
	}

	url, err := EndpointMiscAirportSolarTimeByDate.Replace(map[string]string{
		"codeType":  codeType,
		"code":      code,
		"dateLocal": dateLocal.Format("2006-01-02T15:04"),
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var output *model.SolarStateContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetGreatCircleDistance retrieves the distance and time between two airports - TIER 2
func (c *Client) GetGreatCircleDistance(ctx context.Context, codeType, codeFrom, codeTo string, aircraftName, flightTimeModel *string) (*model.AirportDistanceTimeContract, error) {
	codeType = strings.ToLower(codeType)
	switch codeType {
	case "iata":
		if len(codeFrom) != 3 {
			return nil, fmt.Errorf("invalid IATA code: %s", codeFrom)
		}
		if len(codeTo) != 3 {
			return nil, fmt.Errorf("invalid IATA code: %s", codeTo)
		}
	case "icao":
		if len(codeFrom) != 4 {
			return nil, fmt.Errorf("invalid ICAO code: %s", codeFrom)
		}
		if len(codeTo) != 4 {
			return nil, fmt.Errorf("invalid ICAO code: %s", codeTo)
		}
	default:
		return nil, fmt.Errorf("invalid code type: %s", codeType)
	}

	url, err := EndpointMiscFlightTime.Replace(map[string]string{
		"codeType": codeType,
		"codeFrom": codeFrom,
		"codeTo":   codeTo,
	})
	if err != nil {
		return nil, err
	}

	params, err := model.NewFlightTimeParams(flightTimeModel, aircraftName)
	if err != nil {
		return nil, err
	}
	qv := make(neturl.Values)
	if err := params.EncodeValues("", &qv); err != nil {
		return nil, err
	}

	req, err := c.getRequest(ctx, url+"?"+qv.Encode())
	if err != nil {
		return nil, err
	}

	var output *model.AirportDistanceTimeContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}
