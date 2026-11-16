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
	"context"

	"github.com/aviation-tragic-au/aerodatabox-client/model"
)

type EndpointIndustry int

const (
	EndpointIndustryFAALADDAircraftStatus EndpointIndustry = iota
)

var endpointIndustry = map[EndpointIndustry]string{
	EndpointIndustryFAALADDAircraftStatus: "/industry/faa-ladd/{id}/status",
}

func (e EndpointIndustry) String() string {
	return endpointIndustry[e]
}

// Replace substitutes all {token} placeholders found in the endpoint URL template
// with the corresponding values from the provided map. An error is returned if
// any token present in the template has no matching key in the map.
func (e EndpointIndustry) Replace(tokens map[string]string) (string, error) {
	return replaceTokens(e.String(), tokens)
}

// GetFAALADDAircraftStatus retrieves the status of a FAA LADD aircraft - TIER 3
func (c *Client) GetFAALADDAircraftStatus(ctx context.Context, id string) (*model.FAALADDAircraftStatusContract, error) {
	url, err := EndpointIndustryFAALADDAircraftStatus.Replace(map[string]string{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	req, err := c.getRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var output *model.FAALADDAircraftStatusContract
	if err := c.executeRequest(req, &output); err != nil {
		return nil, err
	}

	return output, nil
}
