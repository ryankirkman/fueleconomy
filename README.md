# api.fueleconomy.io

A modern API on top of [fueleconomy.gov](https://www.fueleconomy.gov/feg/ws/index.shtml) fuel efficiency data.

Get the latest fuel economy data for any car on the U.S. market in clean JSON format.

## Endpoints

### Many Vehicle GET

`GET http://fueleconomy.io/vehicles`

Supported query parameters for multi vehicle get:

**Exact match parameters**

- year

**Fuzzy match parameters**

Queries for these parameters are case insensitive and match with a wild card suffix. So a search for model 'grand cherokee' will return 'Grand Cherokee 2WD' and 'Grand Cherokee 4WD', for example.

- make
- model


#### Response Format

See vehicle object format in Single Vehicle GET response format below.

[http://fueleconomy.io/vehicles?make=jeep&model=grand+cherokee&year=2007](http://fueleconomy.io/vehicles?make=jeep&model=grand+cherokee&year=2007)

```javascript
{
    "vehicles": [
        ...
    ]
}
```

### Single Vehicle GET

`GET http://fueleconomy.io/vehicle/{id}`

#### Response Format

[http://fueleconomy.io/vehicle/23855](http://fueleconomy.io/vehicle/23855)

See `models/vehicle.go` for field descriptions and mapping to fueleconomy.gov fields.

```javascript
{
    "vehicle": {
        "cylinders": 8,
        "driveAxleType": "4-Wheel or All-Wheel Drive",
        "emissionsInfo": [
            {
                "emissionStdCode": "B10",
                "emissionsStdTxt": "BIN 10",
                "engineFamilyId": "7CRXT04.7PSP",
                "f1SmogRating": 1,
                "f2SmogRating": 1,
                "salesArea": 3,
                "updated": "2015-10-02T09:02:36.370078Z"
            }
        ],
        "engDisplacement": 4.7,
        "epaCreatedOn": "2013-01-01T05:00:00Z",
        "epaID": 23855,
        "epaModifiedOn": "2013-01-01T05:00:00Z",
        "fuelType": "Gasoline or E85",
        "fuels": [
            {
                "barrelsPerYear": 21.974,
                "co2Tailpipe": 592.4666666666667,
                "fuelCost": 2350,
                "fuelType": "Regular Gasoline",
                "mpgCity": 14,
                "mpgCityUnadj": 16.7,
                "mpgHighway": 19,
                "mpgHighwayUnadj": 25.6
            },
            {
                "barrelsPerYear": 7.4910000000000005,
                "co2Tailpipe": 620.2,
                "fuelCost": 3200,
                "fuelType": "E85",
                "mpgCity": 9,
                "mpgCityUnadj": 10.7,
                "mpgHighway": 12,
                "mpgHighwayUnadj": 17.2,
                "range": 230
            }
        ],
        "make": "Jeep",
        "model": "Grand Cherokee 4WD",
        "sizeClass": "Sport Utility Vehicle - 4WD",
        "transDscr": "CLKUP",
        "transition": "Automatic 5-spd",
        "updated": "2015-09-29T17:03:23.841714Z",
        "year": 2007
    }
}
```

## Under the hood

Built using go programming language, syncing raw datasets from fueleconomy.gov on a daily basis.

Minimal dependencies:
- [gorilla/mux](https://github.com/gorilla/mux) (excellent router)
- [lib/pq](https://github.com/lib/pq) (postgres driver)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) (sqlite3 driver for testing)
- [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate) (migrations tool)

Custom tooling:
- go struct -> SQL persistence library [(github.com/teasherm/fueleconomy/srm)](https://github.com/teasherm/fueleconomy/srm)
- worker pool and task queue [(github.com/teasherm/fueleconomy/workers)](https://github.com/teasherm/fueleconomy/workers)
