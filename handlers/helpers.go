package handlers

import (
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/models"
)

func checkErr(err error, w http.ResponseWriter) {
	if err != nil {
		global.Logger.Println("Error: ", err)
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func getIntFromQueryVals(queryVals url.Values, param string) int {
	value := queryVals.Get(param)
	if parsed, err := strconv.Atoi(value); err != nil {
		return 0
	} else {
		return parsed
	}
}

func getProfileFromQueryVals(queryVals url.Values) models.DrivingProfile {
	profile := models.DrivingProfile{
		CityShare:    models.CityShareDefault,
		HighwayShare: models.HighwayShareDefault,
		MilesPerYear: models.MilesPerYearDefault,
	}

	cityShare := getIntFromQueryVals(queryVals, "cityShare")
	if cityShare > 0 {
		profile.CityShare = cityShare
	}

	highwayShare := getIntFromQueryVals(queryVals, "highwayShare")
	if highwayShare > 0 {
		profile.HighwayShare = highwayShare
	}

	milesPerYear := getIntFromQueryVals(queryVals, "milesPerYear")
	if milesPerYear > 0 {
		profile.MilesPerYear = milesPerYear
	}

	return profile
}

func getMostRecentFuelPrices() (fp models.FuelPrices) {
	query := "SELECT * FROM fuel_prices WHERE updated = (SELECT MAX(updated) from fuel_prices)"
	global.Db.SelectOne(&fp, query)
	return fp
}

func mapFromQueryVals(queryString url.Values, params []string) (res map[string]string) {
	res = make(map[string]string)
	for _, q := range params {
		val := queryString.Get(q)
		if val != "" {
			res[q] = val
		}
	}
	return res
}

func maxInt(first int, second int) int {
	return int(math.Max(float64(first), float64(second)))
}
