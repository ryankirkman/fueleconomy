# fueleconomy.io

A modern API on top of [fueleconomy.gov](https://www.fueleconomy.gov/feg/ws/index.shtml) fuel efficiency data.

Get the latest fuel economy data for any car on the U.S. market in clean JSON format.

## Endpoints

### Many Vehicle GET

`GET http://fueleconomy.io/vehicles`

Supported query parameters for multi vehicle get:

**Search parameters**

Queries for "fuzzy" parameters are case insensitive and match with a wild card suffix. So a search for model 'grand cherokee' will return 'Grand Cherokee 2WD' and 'Grand Cherokee 4WD', for example.

- make (fuzzy)
- model (fuzzy)
- year

**Driving profile parameters**

These parameters affect the calculation of combined mpg and annual fuel cost

- cityShare - Percentage city driving (Default: 55%)
- highwayShare - Percentage highway driving (Default: 45%)
- milesPerYear - Miles driven per year (Default: 15,000)

**Pagination parameters**

- page - Page number (Default: 1)
- pageLength - Number of results per page (Default: 10, Max: 100)


#### Response Format

[http://fueleconomy.io/vehicles?make=jeep&page=2&year=2016](http://fueleconomy.io/vehicles?make=jeep&model=grand+cherokee&year=2007)

```javascript
{
    "meta": {
        "nextPage": "http://fueleconomy.io/vehicles?make=Jeep&page=3&year=2016",
        "page": 2,
        "pageLength": 10,
        "prevPage": "http://fueleconomy.io/vehicles?make=Jeep&page=1&year=2016",
        "totalPages": 3,
        "totalResults": 28
    },
    "profile": {
        "cityShare": 55,
        "highwayShare": 45,
        "milesPerYear": 15000
    },
    "vehicles": [
        {
            "cylinders": 4,
            "driveAxleType": "4-Wheel Drive",
            "emissionsInfo": [
                {
                    "emissionStdCode": "B4",
                    "emissionsStdTxt": "Bin 4",
                    "engineFamilyId": "GCRXT02.45P2",
                    "f1SmogRating": 6,
                    "salesArea": 7,
                    "updated": "2015-10-18T22:17:00.826353Z"
                },
                {
                    "emissionStdCode": "B4",
                    "emissionsStdTxt": "Bin 4",
                    "engineFamilyId": "GCRXT02.45P2",
                    "f1SmogRating": 6,
                    "salesArea": 3,
                    "updated": "2015-10-18T22:17:00.825866Z"
                }
            ],
            "engDisplacement": 2.4,
            "engDscr": "SIDI",
            "engID": 505,
            "epaCreatedOn": "2015-06-24T04:00:00Z",
            "epaID": 36431,
            "epaModifiedOn": "2015-08-18T04:00:00Z",
            "fuelEconomyScore": 5,
            "fuelType": "Regular",
            "fuels": [
                {
                    "barrelsPerYear": 14.96,
                    "co2": 398,
                    "co2Tailpipe": 398,
                    "fuelCost": 1539,
                    "fuelType": "Regular Gasoline",
                    "ghgScore": 5,
                    "mpgCity": 20,
                    "mpgCityUnrounded": 20.1804,
                    "mpgComb": 22.7,
                    "mpgHighway": 26,
                    "mpgHighwayUnrounded": 25.7164
                }
            ],
            "make": "Jeep",
            "manufacturerCode": "CRX",
            "model": "Patriot 4WD",
            "sizeClass": "Small Sport Utility Vehicle 4WD",
            "transition": "Automatic 6-spd",
            "updated": "2015-09-19T22:33:38.823372Z",
            "year": 2016
        },
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
    "profile": {
        "cityShare": 55,
        "highwayShare": 45,
        "milesPerYear": 15000
    },
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
                "updated": "2015-10-18T22:16:50.099966Z"
            }
        ],
        "engDisplacement": 4.7,
        "epaCreatedOn": "2013-01-01T05:00:00Z",
        "epaID": 23855,
        "epaModifiedOn": "2013-01-01T05:00:00Z",
        "fuelType": "Gasoline or E85",
        "fuels": [
            {
                "barrelsPerYear": 21.97,
                "co2Tailpipe": 592.47,
                "fuelCost": 2150,
                "fuelType": "Regular Gasoline",
                "mpgCity": 14,
                "mpgComb": 16.25,
                "mpgHighway": 19
            },
            {
                "barrelsPerYear": 7.49,
                "co2Tailpipe": 620.2,
                "fuelCost": 3086,
                "fuelType": "E85",
                "mpgCity": 9,
                "mpgComb": 10.35,
                "mpgHighway": 12
            }
        ],
        "make": "Jeep",
        "model": "Grand Cherokee 4WD",
        "sizeClass": "Sport Utility Vehicle - 4WD",
        "transDscr": "CLKUP",
        "transition": "Automatic 5-spd",
        "updated": "2015-09-19T22:33:15.043763Z",
        "year": 2007
    }
}
```

## Under the hood

Syncs raw datasets from fueleconomy.gov on a daily basis.

Minimal dependencies:
- [gorilla/mux](https://github.com/gorilla/mux) (excellent router)
- [lib/pq](https://github.com/lib/pq) (postgres driver)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) (sqlite3 driver for testing)
- [rubenv/sql-migrate](https://github.com/rubenv/sql-migrate) (migrations tool)

Custom tooling:
- go struct -> SQL persistence library [(github.com/teasherm/fueleconomy/srm)](https://github.com/teasherm/fueleconomy/tree/master/srm)
- worker pool and task queue [(github.com/teasherm/fueleconomy/workers)](https://github.com/teasherm/fueleconomy/tree/master/workers)
