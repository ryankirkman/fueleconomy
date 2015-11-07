package handlers

import (
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/models"
)

// Error check helper

func checkErr(err error, w http.ResponseWriter) {
	if err != nil {
		global.Logger.Println("Error: ", err)
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
}

// Search param data struct and parser

type searchParam struct {
	name      string
	converter func(string) (interface{}, error)
}

func intConverter(in string) (interface{}, error) {
	return strconv.Atoi(in)
}

func extractSearchParams(queryVals url.Values, params []searchParam) map[string]interface{} {
	out := make(map[string]interface{})
	for _, param := range params {
		val := queryVals.Get(param.name)
		if val != "" {
			if conv, err := param.converter(val); err != nil {
				out[param.name] = conv
			}
		}
	}
	return out
}

func extractStringParams(queryVals url.Values, params []string) map[string]string {
	out := make(map[string]string)
	for _, param := range params {
		val := queryVals.Get(param)
		if val != "" {
			out[param] = val
		}
	}
	return out
}

// Driving profile parser

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

// Fuel prices retriever

func getMostRecentFuelPrices() (fp models.FuelPrices) {
	query := "SELECT * FROM fuel_prices WHERE updated = (SELECT MAX(updated) from fuel_prices)"
	global.Db.SelectOne(&fp, query)
	return fp
}

// Get maximum of two integers

func maxInt(first int, second int) int {
	return int(math.Max(float64(first), float64(second)))
}
