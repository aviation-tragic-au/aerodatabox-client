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
	"time"

	"github.com/deanstalker-stream/aerodatabox-client/model"
)

// EndpointStatistical represents the statistical endpoints.
type EndpointStatistical int

const (
	EndpointStatisticalAirportDelays EndpointStatistical = iota
	EndpointStatisticalAirportDelaysByDate
	EndpointStatisticalAirportDelaysByDateRange
	EndpointStatisticalAirportRouteStatistics
	EndpointStatisticalAirportRouteStatisticsByDate
	EndpointStatisticalGlobalDelaysCurrent
	EndpointStatisticalGlobalDelaysByDate
	EndpointStatisticalFlightDelayStatistics
)

var endpointStatistical = map[EndpointStatistical]string{
	EndpointStatisticalAirportDelays:                "/airports/{codeType}/{code}/delays",
	EndpointStatisticalAirportDelaysByDate:          "/airports/{codeType}/{code}/delays/{dateLocal}",
	EndpointStatisticalAirportDelaysByDateRange:     "/airports/{codeType}/{code}/delays/{dateFrom}/{dateTo}",
	EndpointStatisticalAirportRouteStatistics:       "/airports/{codeType}/{code}/stats/routes/daily",
	EndpointStatisticalAirportRouteStatisticsByDate: "/airports/{codeType}/{code}/stats/routes/daily/{dateLocal}",
	EndpointStatisticalGlobalDelaysCurrent:          "/airports/delays",
	EndpointStatisticalGlobalDelaysByDate:           "/airports/delays/{dateUtc}",
	EndpointStatisticalFlightDelayStatistics:        "/flights/{number}/delays",
}

// String returns the endpoint URL.
func (e EndpointStatistical) String() string {
	return endpointStatistical[e]
}

// Replace substitutes all {token} placeholders found in the endpoint URL template
// with the corresponding values from the provided map. An error is returned if
// any of the tokens are not found in the map.
func (e EndpointStatistical) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// GetAirportDelays retrieves the delays of an airport - TIER 3
func (c *Client) GetAirportDelays(code, codeType string, dateLocal *time.Time) (*model.AirportDelayContract, error) {
	var url string
	var err error
	if dateLocal != nil {
		url, err = EndpointStatisticalAirportDelaysByDate.Replace(map[string]string{
			"codeType":  codeType,
			"code":      code,
			"dateLocal": dateLocal.Format("2006-01-02T15:04:05"),
		})
	} else {
		url, err = EndpointStatisticalAirportDelays.Replace(map[string]string{
			"codeType": codeType,
			"code":     code,
		})
	}
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}

	var output *model.AirportDelayContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetAirportDelaysInRange retrieves the delays of an airport within a given date range - TIER 3
func (c *Client) GetAirportDelaysInRange(code, codeType string, dateFromLocal, dateToLocal time.Time) ([]*model.AirportDelayContract, error) {
	url, err := EndpointStatisticalAirportDelaysByDateRange.Replace(map[string]string{
		"codeType": codeType,
		"code":     code,
		"dateFrom": dateFromLocal.Format("2006-01-02T15:04:05"),
		"dateTo":   dateToLocal.Format("2006-01-02T15:04:05"),
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var output []*model.AirportDelayContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetAirportRoutes retrieves the daily route statistics of an airport - TIER 3
func (c *Client) GetAirportRoutes(code, codeType string) (*model.DailyRouteStatContract, error) {
	url, err := EndpointStatisticalAirportRouteStatistics.Replace(map[string]string{
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
	var output *model.DailyRouteStatContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetAirportRoutesByDate retrieves the daily route statistics of an airport for a specific date - TIER 3
func (c *Client) GetAirportRoutesByDate(code, codeType string, dateLocal time.Time) (*model.DailyRouteStatContract, error) {
	url, err := EndpointStatisticalAirportRouteStatisticsByDate.Replace(map[string]string{
		"codeType":  codeType,
		"code":      code,
		"dateLocal": dateLocal.Format("2006-01-02T15:04:05"),
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var output *model.DailyRouteStatContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetGlobalDelays retrieves the delays of an airport - TIER 3
func (c *Client) GetGlobalDelays() (*model.AirportDelayContract, error) {
	url := EndpointStatisticalGlobalDelaysCurrent.String()
	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var output *model.AirportDelayContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetGlobalDelaysByDate retrieves the delays of an airport for a specific date - TIER 3
func (c *Client) GetGlobalDelaysByDate(dateUtc time.Time) (*model.AirportDelayContract, error) {
	url, err := EndpointStatisticalGlobalDelaysByDate.Replace(map[string]string{
		"dateUtc": dateUtc.Format("2006-01-02T15:04:05"),
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var output *model.AirportDelayContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// GetFlightDelayStatistics retrieves the flight delay statistics of a flight - TIER 3
func (c *Client) GetFlightDelayStatistics(number string) (*model.FlightDelayContract, error) {
	url, err := EndpointStatisticalFlightDelayStatistics.Replace(map[string]string{
		"number": number,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(url)
	if err != nil {
		return nil, err
	}
	var output *model.FlightDelayContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}
	return output, nil
}
