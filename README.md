# aerodatabox-client

[![Test](https://github.com/deanstalker-stream/aerodatabox-client/actions/workflows/test.yml/badge.svg)](https://github.com/deanstalker-stream/aerodatabox-client/actions/workflows/test.yml) [![Lint](https://github.com/deanstalker-stream/aerodatabox-client/actions/workflows/lint.yml/badge.svg)](https://github.com/deanstalker-stream/aerodatabox-client/actions/workflows/lint.yml)

An unofficial, community-developed Go client for the [AeroDataBox API](https://www.aerodatabox.com/). Not affiliated with, endorsed by, or sponsored by AeroDataBox or Timber Bytes LLP. Users are responsible for complying with the AeroDataBox Terms of Service.

## Installation

```bash
go get github.com/deanstalker-stream/aerodatabox-client
```

## Configuration

The client reads two values from `model.Config`:

| Field  | Description                          |
|--------|--------------------------------------|
| `Host` | RapidAPI host (e.g. `aerodatabox.p.rapidapi.com`) |
| `Key`  | Your RapidAPI subscription key       |

If you use [Viper](https://github.com/spf13/viper) for configuration, call `client.BindEnvs(v)` once at startup to register the environment variable mappings:

| Environment variable      | Config key               |
|---------------------------|--------------------------|
| `FEED_AERODATABOX_HOST`   | `feed.aerodatabox.host`  |
| `FEED_AERODATABOX_KEY`    | `feed.aerodatabox.key`   |

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    aero "github.com/deanstalker-stream/aerodatabox-client/client"
    "github.com/deanstalker-stream/aerodatabox-client/model"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()

    cfg := &model.Config{
        Host: "aerodatabox.p.rapidapi.com",
        Key:  "YOUR_RAPIDAPI_KEY",
    }

    c, err := aero.NewClient(logger, cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Look up an airport by ICAO code
    airport, err := c.GetAirportByCode("EGLL", "icao", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(airport.FullName)

    // Get upcoming departures from Heathrow (next 2 hours)
    flights, err := c.GetFlightInformationDepartureByRelativeTime(
        "icao", "EGLL",
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

### Using Viper

```go
import (
    "github.com/spf13/viper"
    aero "github.com/deanstalker-stream/aerodatabox-client/client"
    "github.com/deanstalker-stream/aerodatabox-client/model"
)

func setup(v *viper.Viper) {
    aero.BindEnvs(v)
    // v will now read Host from FEED_AERODATABOX_HOST and Key from FEED_AERODATABOX_KEY

    var cfg model.Config
    _ = v.UnmarshalKey("feed.aerodatabox", &cfg)
}
```

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
| `GetAirportByCode(code, codeType, withTime, withRunways)` | Airport by IATA or ICAO code | 1 |
| `GetAirportRunways(code, codeType)` | Runway details for an airport | 1 |
| `GetAirportByLocation(lat, lon, limit, radiusKm, withFlightInfoOnly)` | Airports within a radius | 2 |
| `GetAirportByIP(q, limit, radiusKM, withFlightInfoOnly)` | Airports near an IP address | 2 |
| `GetAirportByTerm(q, limit, withFlightInfoOnly, withSearchByCode)` | Airport search by term | 2 |

### Flights

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetFlightStatus(searchBy, searchParameter, dateLocalRole, ...)` | Current status of a flight | 2 |
| `GetFlightStatusByDate(searchBy, searchParameter, dateLocal, ...)` | Flight status on a specific date | 2 |
| `GetFlightScheduleHistory(searchBy, searchParameter, dateFrom, dateTo, ...)` | Schedule history for a flight | 2 |
| `GetFlightDepartureDates(searchBy, searchParameter)` | All departure dates for a flight | 2 |
| `GetFlightDepartureDateRange(searchBy, searchParameter, from, to)` | Departure dates in a range | 2 |
| `GetFlightInformationDepartureByLocalTime(codeType, code, from, to, ...)` | Airport FIDS by local time window | 2 |
| `GetFlightInformationDepartureByRelativeTime(codeType, code, ...)` | Airport FIDS relative to now | 2 |
| `SearchFlightByTerm(term, limit)` | Search flights by term | 2 |

`searchBy` accepts: `"Number"`, `"Reg"`, `"CallSign"`, `"Icao24"` (see `model.FlightSearchByEnum`).

### Aircraft

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetSingleAircraft(searchBy, searchParameter, withRegistrations, withImage)` | Single aircraft record | 1 |
| `GetAircraft(searchBy, searchParameter, withRegistrations, withImage)` | All aircraft matching criteria | — |
| `GetAircraftHistory(searchBy, searchParameter)` | Registration history for an aircraft | — |
| `GetAircraftByTerm(term, limit)` | Aircraft search by term | — |
| `GetAirlineFleetByCode(code, pageSize, pageOffset, withRegistrations)` | Fleet for an airline | — |

`searchBy` accepts: `"Id"`, `"Reg"`, `"Icao24"` (see `model.AircraftSearchByEnum`).

### Statistics

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetAirportDelays(code, codeType, dateLocal)` | Delay statistics for an airport | 3 |
| `GetAirportDelaysInRange(code, codeType, dateFrom, dateTo)` | Delay statistics over a date range | 3 |
| `GetAirportRoutes(code, codeType)` | Active routes for an airport | 3 |
| `GetAirportRoutesByDate(code, codeType, dateLocal)` | Routes for an airport on a date | 3 |
| `GetGlobalDelays()` | Current global delay statistics | 3 |
| `GetGlobalDelaysByDate(dateUtc)` | Global delays on a given date | 3 |
| `GetFlightDelayStatistics(number)` | Historical delay stats for a flight number | 3 |

### Miscellaneous

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetAirportCurrentTime(code, codeType)` | Current local time at an airport | 1 |
| `GetAirportSolarTime(code, codeType)` | Sunrise/sunset at an airport | 1 |
| `GetAirportSolarTimeByDate(code, codeType, dateLocal)` | Solar times on a specific date | 1 |
| `GetGreatCircleDistance(codeType, from, to, aircraftName, flightTimeModel)` | Great-circle distance between airports | 2 |

### Health

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetDataFeedStatus(service, withHttpCode)` | Status of a named data feed | Free |
| `GetDataFeedStatusByICAO(icao)` | Feed status for a specific airport | Free |
| `GetDataFeedSupportedAirports(service)` | Airports supported by a feed | Free |

### Flight Alerts (Webhooks)

| Method | Description |
|--------|-------------|
| `CreateWebhookSubscription(subjectID, subjectType, callbackURL, maxRetries)` | Subscribe to flight events |
| `GetWebhookSubscription(subscriptionID)` | Get a subscription |
| `DeleteWebhookSubscription(subscriptionID)` | Cancel a subscription |
| `ListWebhookSubscriptions()` | List all subscriptions |
| `GetWebhookSubscriptionsBalance()` | Check credit balance |
| `RefillWebhookSubscriptionsBalance(credits)` | Add credits |

### Industry Data

| Method | Description | Tier |
|--------|-------------| ---- |
| `GetFAALADDAircraftStatus(id)` | FAA LADD aircraft status | 3 |

## Error Handling

All methods return a `*model.APIError` for non-2xx responses. It includes the HTTP status code and the message from the API:

```go
flights, err := c.GetFlightStatus("Number", "BA123", model.FlightDirectionArrival, false, false, false)
if err != nil {
    var apiErr *model.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API error %d: %s\n", apiErr.StatusCode, apiErr.Message)
    }
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
