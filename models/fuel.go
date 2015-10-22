package models

type Fuel struct {
	BarrelsPerYear      float64 `json:"barrelsPerYear,omitempty"`      // annual petroleum consumption in barrels for fuelType1 (1)
	Co2                 float64 `json:"co2,omitempty"`                 // tailpipe CO2 in grams/mile for fuelType1 (5)
	Co2Tailpipe         float64 `json:"co2Tailpipe,omitempty"`         // tailpipe CO2 in grams/mile for fuelType1 (5)
	ECity               float64 `json:"eCity,omitempty"`               // city electricity consumption in kw-hrs/100 miles
	EComb               float64 `json:"eComb,omitempty"`               // combined electricity consumption in kw-hrs/100 miles
	EHighway            float64 `json:"eHighway,omitempty"`            // highway electricity consumption in kw-hrs/100 miles
	FuelCost            int     `json:"fuelCost,omitempty"`            // annual fuel cost for fuelType1 ($) (7)
	FuelType            string  `json:"fuelType,omitempty"`            // fuel type 1. For single fuel vehicles, this will be the only fuel. For dual fuel vehicles, this will be the conventional fuel.
	GhgScore            int     `json:"ghgScore,omitempty"`            // EPA GHG score (-1 = Not available)
	MpgCity             float64 `json:"mpgCity,omitempty"`             // city MPG for fuelType1 (2)
	MpgCityUnadj        float64 `json:"-"`                             // unadjusted city MPG for fuelType1; see the description of the EPA test procedures
	MpgCityUnrounded    float64 `json:"mpgCityUnrounded,omitempty"`    // unrounded city MPG for fuelType1 (2), (3)
	MpgComb             float64 `json:"mpgComb, omitempty"`            // combined MPG for fuelType1 (2)
	MpgCombUnrounded    float64 `json:"-"`                             // unrounded combined MPG for fuelType1 (2), (3)
	MpgHighway          float64 `json:"mpgHighway,omitempty"`          // highway MPG for fuelType1 (2)
	MpgHighwayUnadj     float64 `json:"-"`                             // unadjusted highway MPG for fuelType1; see the description of the EPA test procedures
	MpgHighwayUnrounded float64 `json:"mpgHighwayUnrounded,omitempty"` // unrounded highway MPG for fuelType1 (2), (3)
	Range               float64 `json:"range,omitempty"`               // EPA range for fuelType2
	RangeCity           float64 `json:"rangeCity,omitempty"`           // EPA city range for fuelType2
	RangeHighway        float64 `json:"rangeHighway,omitempty"`        // EPA highway range for fuelType2
	PhevCDCity          float64 `json:"phevCDCity,omitempty"`          // city gasoline consumption (gallons/100 miles) in charge depleting mode (4)
	PhevCDComb          float64 `json:"phevCDComb,omitempty"`          // combined gasoline consumption (gallons/100 miles) in charge depleting mode (4)
	PhevCDHighway       float64 `json:"phevCDHighway,omitempty"`       // highway gasoline consumption (gallons/100miles) in charge depleting mode (4)
	PhevMpgCity         float64 `json:"phevMpgCity,omitempty"`         // EPA composite gasoline-electricity city MPGe for plug-in hybrid vehicles
	PhevMpgComb         float64 `json:"phevMpgComb,omitempty"`         // EPA composite gasoline-electricity combined city-highway MPGe for plug-in hybrid vehicles
	PhevMpgHighway      float64 `json:"phevMpgHighway,omitempty"`      // EPA composite gasoline-electricity highway MPGe for plug-in hybrid vehicles
	PhevUFCity          float64 `json:"phevUFCity,omitempty"`          // EPA city utility factor (share of electricity) for PHEV
	PhevUFComb          float64 `json:"phevUFComb,omitempty"`          // EPA combined utility factor (share of electricity) for PHEV
	PhevUFHighway       float64 `json:"phevUFHighway,omitempty"`       // EPA highway utility factor (share of electricity) for PHEV
}

func CalculateFuelData(v *Vehicle, d DrivingProfile, fp FuelPrices) (fuels []Fuel) {
	fuels = append(fuels, calculateFuelOne(v, d, fp))
	if v.F2FuelType != "" {
		fuels = append(fuels, calculateFuelTwo(v, d, fp))
	}
	return fuels
}

func calculateFuelOne(v *Vehicle, d DrivingProfile, fp FuelPrices) (fuel Fuel) {
	fuel = Fuel{
		Co2:                 toFixed(v.F1Co2, 2),
		Co2Tailpipe:         toFixed(v.F1Co2Tailpipe, 2),
		ECity:               v.ECity,
		EHighway:            v.EHighway,
		FuelType:            v.F1FuelType,
		GhgScore:            v.F1GhgScore,
		MpgCity:             v.F1MpgCity,
		MpgCityUnadj:        v.F1MpgCityUnadj,
		MpgCityUnrounded:    v.F1MpgCityUnrounded,
		MpgHighway:          v.F1MpgHighway,
		MpgHighwayUnadj:     v.F1MpgHighwayUnadj,
		MpgHighwayUnrounded: v.F1MpgHighwayUnrounded,
		Range:               v.F1Range,
		PhevCDCity:          v.PhevCDCity,
		PhevCDHighway:       v.PhevCDHighway,
		PhevMpgCity:         v.PhevMpgCity,
		PhevMpgHighway:      v.PhevMpgHighway,
		PhevUFCity:          v.PhevUFCity,
		PhevUFHighway:       v.PhevUFHighway,
	}

	if v.F1BarrelsPerYear > 0.0 {
		fuel.BarrelsPerYear = calculateBarrelsPerYear(v.F1BarrelsPerYear, d.MilesPerYear)
	}

	if v.EComb > 0.0 {
		fuel.EComb = calculateEComb(v.ECity, d.CityShare, v.EHighway, d.HighwayShare)
	}

	if v.F1MpgComb > 0.0 {
		fuel.MpgComb = calculateMpgComb(v.F1MpgCity, d.CityShare, v.F1MpgHighway,
			d.HighwayShare)
	}

	if v.PhevCDComb > 0.0 {
		fuel.PhevCDComb = calculateMpgComb(v.PhevCDCity, d.CityShare, v.PhevCDHighway,
			d.HighwayShare)
	}

	if v.PhevMpgComb > 0.0 {
		fuel.PhevMpgComb = calculateMpgComb(v.PhevMpgCity, d.CityShare, v.PhevMpgHighway,
			d.HighwayShare)
	}

	if v.PhevUFComb > 0.0 {
		fuel.PhevUFComb = calculateMpgComb(v.PhevUFCity, d.CityShare, v.PhevUFHighway,
			d.HighwayShare)
	}

	fuelPrice := fuelPriceFromName(v.F1FuelType, fp)
	if fuelPrice > 0.0 {
		switch v.F1FuelType {
		default:
			fuel.FuelCost = calculateFuelCost(fuelPrice, d.MilesPerYear, fuel.MpgComb)
		case "Electricity":
			fuel.FuelCost = calculateFuelCostElectricity(fuelPrice, d.MilesPerYear, fuel.EComb)
		}
	}

	return fuel
}

func calculateFuelTwo(v *Vehicle, d DrivingProfile, fp FuelPrices) (fuel Fuel) {
	fuel = Fuel{
		Co2:                 toFixed(v.F2Co2, 2),
		Co2Tailpipe:         toFixed(v.F2Co2Tailpipe, 2),
		FuelType:            v.F2FuelType,
		GhgScore:            v.F2GhgScore,
		MpgCity:             v.F2MpgCity,
		MpgCityUnadj:        v.F2MpgCityUnadj,
		MpgCityUnrounded:    v.F2MpgCityUnrounded,
		MpgHighway:          v.F2MpgHighway,
		MpgHighwayUnadj:     v.F2MpgHighwayUnadj,
		MpgHighwayUnrounded: v.F2MpgHighwayUnrounded,
		RangeCity:           v.F2RangeCity,
		RangeHighway:        v.F2RangeHighway,
	}

	if v.F2BarrelsPerYear > 0.0 {
		fuel.BarrelsPerYear = calculateBarrelsPerYear(v.F2BarrelsPerYear, d.MilesPerYear)
	}

	if v.F2MpgComb > 0.0 {
		fuel.MpgComb = calculateMpgComb(v.F2MpgCity, d.CityShare, v.F2MpgHighway,
			d.HighwayShare)
	}

	if v.F2Range > 0.0 {
		fuel.Range = calculateRange(v.F2RangeCity, d.CityShare, v.F2RangeHighway, d.HighwayShare)
	}

	fuelPrice := fuelPriceFromName(v.F2FuelType, fp)
	if fuelPrice > 0.0 {
		switch v.F2FuelType {
		default:
			fuel.FuelCost = calculateFuelCost(fuelPrice, d.MilesPerYear, fuel.MpgComb)
		case "Electricity":
			fuel.FuelCost = calculateFuelCostElectricity(fuelPrice, d.MilesPerYear, fuel.EComb)
		}
	}

	return fuel
}
