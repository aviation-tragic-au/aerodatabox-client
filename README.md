# aerodatabox-client

[![Test](https://github.com/aviation-tragic-au/aerodatabox-client/actions/workflows/test.yml/badge.svg)](https://github.com/aviation-tragic-au/aerodatabox-client/actions/workflows/test.yml) [![Lint](https://github.com/aviation-tragic-au/aerodatabox-client/actions/workflows/lint.yml/badge.svg)](https://github.com/aviation-tragic-au/aerodatabox-client/actions/workflows/lint.yml)

An unofficial, community-developed Go client for the [AeroDataBox API](https://www.aerodatabox.com/). Not affiliated with, endorsed by, or sponsored by AeroDataBox or Timber Bytes LLP. Users are responsible for complying with the AeroDataBox Terms of Service.

## Installation

```bash
go get github.com/aviation-tragic-au/aerodatabox-client
```

## Configuration

`client.Config` holds everything. Only a key is required:

| Field        | Description |
|--------------|-------------|
| `Key`        | Your RapidAPI key |
| `KeyFunc`    | Called on every request in place of `Key`, to pick up a key changed at runtime |
| `Host`       | RapidAPI host, sent as `X-RapidAPI-Host` (default `aerodatabox.p.rapidapi.com`) |
| `URL`        | API base (default `https://<Host>`); set it for a proxy or a test server |
| `HTTPClient` | Your own `*http.Client` (the default times out after 10 seconds) |

`Config` carries `json` and `mapstructure` tags, so it can be decoded
straight from a config file or [Viper](https://github.com/spf13/viper)
(`v.UnmarshalKey("aerodatabox", &cfg)`); the client itself doesn't depend
on either.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    aero "github.com/aviation-tragic-au/aerodatabox-client/client"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()

    c, err := aero.NewClient(logger, aero.Config{Key: "YOUR_RAPIDAPI_KEY"})
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Look up an airport by ICAO code
    airport, err := c.GetAirportByCode(ctx, "EGLL", "icao", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(airport.FullName)

    // Get upcoming departures from Heathrow (next 2 hours)
    flights, err := c.GetFlightInformationDepartureByRelativeTime(
        ctx, "icao", "EGLL",
        nil, nil, nil, nil, nil, nil, nil,
        nil, nil,
    )
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range flights.Departures {
        fmt.Println(f.Number, f.Status.String())
    }
}
```

Every method takes a `context.Context` first, for cancellation and
deadlines.

## API Reference

Endpoints are billed in units based on their computational complexity. See [AeroDataBox pricing](https://aerodatabox.com/pricing/) for plan unit allocations.

| Tier | Units per request |
|------|-------------------|
| Free | 0 |
| 1    | 1 |
| 2    | 2 |
| 3    | 6 |

### Airports

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetAirportByCode(ctx, code, codeType, withTime, withRunways)` | Airport by IATA or ICAO code | 1 |
| `GetAirportRunways(ctx, code, codeType)` | Runway details for an airport | 1 |
| `GetAirportByLocation(ctx, lat, lon, limit, radiusKm, withFlightInfoOnly)` | Airports within a radius | 2 |
| `GetAirportByIP(ctx, q, limit, radiusKM, withFlightInfoOnly)` | Airports near an IP address | 2 |
| `GetAirportByTerm(ctx, q, limit, withFlightInfoOnly, withSearchByCode)` | Airport search by term | 2 |

### Flights

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetFlightStatus(ctx, searchBy, searchParameter, dateLocalRole, ...)` | Current status of a flight | 2 |
| `GetFlightStatusByDate(ctx, searchBy, searchParameter, dateLocal, ...)` | Flight status on a specific date | 2 |
| `GetFlightScheduleHistory(ctx, searchBy, searchParameter, dateFrom, dateTo, ...)` | Schedule history for a flight | 2 |
| `GetFlightDepartureDates(ctx, searchBy, searchParameter)` | All departure dates for a flight | 2 |
| `GetFlightDepartureDateRange(ctx, searchBy, searchParameter, from, to)` | Departure dates in a range | 2 |
| `GetFlightInformationDepartureByLocalTime(ctx, codeType, code, from, to, ...)` | Airport FIDS by local time window | 2 |
| `GetFlightInformationDepartureByRelativeTime(ctx, codeType, code, ...)` | Airport FIDS relative to now | 2 |
| `SearchFlightByTerm(ctx, term, limit)` | Search flights by term | 2 |

`searchBy` accepts: `"Number"`, `"Reg"`, `"CallSign"`, `"Icao24"` (see `model.FlightSearchByEnum`).

### Aircraft

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetSingleAircraft(ctx, searchBy, searchParameter, withRegistrations, withImage)` | Single aircraft record | 1 |
| `GetAircraft(ctx, searchBy, searchParameter, withRegistrations, withImage)` | All aircraft matching criteria | — |
| `GetAircraftHistory(ctx, searchBy, searchParameter)` | Registration history for an aircraft | — |
| `GetAircraftByTerm(ctx, term, limit)` | Aircraft search by term | — |
| `GetAirlineFleetByCode(ctx, code, pageSize, pageOffset, withRegistrations)` | Fleet for an airline | — |

`searchBy` accepts: `"Id"`, `"Reg"`, `"Icao24"` (see `model.AircraftSearchByEnum`).

### Statistics

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetAirportDelays(ctx, code, codeType, dateLocal)` | Delay statistics for an airport | 3 |
| `GetAirportDelaysInRange(ctx, code, codeType, dateFrom, dateTo)` | Delay statistics over a date range | 3 |
| `GetAirportRoutes(ctx, code, codeType)` | Active routes for an airport | 3 |
| `GetAirportRoutesByDate(ctx, code, codeType, dateLocal)` | Routes for an airport on a date | 3 |
| `GetGlobalDelays(ctx)` | Current global delay statistics | 3 |
| `GetGlobalDelaysByDate(ctx, dateUtc)` | Global delays on a given date | 3 |
| `GetFlightDelayStatistics(ctx, number)` | Historical delay stats for a flight number | 3 |

### Miscellaneous

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetAirportCurrentTime(ctx, code, codeType)` | Current local time at an airport | 1 |
| `GetAirportSolarTime(ctx, code, codeType)` | Sunrise/sunset at an airport | 1 |
| `GetAirportSolarTimeByDate(ctx, code, codeType, dateLocal)` | Solar times on a specific date | 1 |
| `GetGreatCircleDistance(ctx, codeType, from, to, aircraftName, flightTimeModel)` | Great-circle distance between airports | 2 |

### Health

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetDataFeedStatus(ctx, service, withHttpCode)` | Status of a named data feed | Free |
| `GetDataFeedStatusByICAO(ctx, icao)` | Feed status for a specific airport | Free |
| `GetDataFeedSupportedAirports(ctx, service)` | Airports supported by a feed | Free |

### Flight Alerts (Webhooks)

| Method | Description |
|--------|-------------|
| `CreateWebhookSubscription(ctx, subjectID, subjectType, callbackURL, maxRetries)` | Subscribe to flight events |
| `GetWebhookSubscription(ctx, subscriptionID)` | Get a subscription |
| `DeleteWebhookSubscription(ctx, subscriptionID)` | Cancel a subscription |
| `ListWebhookSubscriptions(ctx)` | List all subscriptions |
| `GetWebhookSubscriptionsBalance(ctx)` | Check credit balance |
| `RefillWebhookSubscriptionsBalance(ctx, credits)` | Add credits |

### Industry Data

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetFAALADDAircraftStatus(ctx, id)` | FAA LADD aircraft status | 3 |

## Error Handling

AeroDataBox doesn't cover every airport or flight. When it has nothing for
a request (a 204, an empty body, or a 404) the method returns an error
matching `aero.ErrNoData`. Any other non-2xx is a `*model.APIError` with
the HTTP status and the API's message (a 404 is both):

```go
flights, err := c.GetFlightStatus(ctx, "Number", "BA123", model.FlightDirectionArrival, false, false, false)
var apiErr *model.APIError
switch {
case errors.Is(err, aero.ErrNoData):
    fmt.Println("AeroDataBox has nothing for BA123")
case errors.As(err, &apiErr):
    fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
case err != nil:
    return err
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

Copyright 2026 Dean Stalker

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
