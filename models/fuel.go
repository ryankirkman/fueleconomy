package models

func ConstructFuelData(v *Vehicle) (fuels []Fuel) {
	fuels = append(fuels, buildFuel1(v))
	if v.F2FuelType != "" {
		fuels = append(fuels, buildFuel2(v))
	}
	return fuels
}

func buildFuel1(v *Vehicle) (fuel Fuel) {
	fuel = Fuel{
		BarrelsPerYear:      v.F1BarrelsPerYear,
		Co2:                 v.F1Co2,
		Co2Tailpipe:         v.F1Co2Tailpipe,
		ECity:               v.ECity,
		EComb:               v.EComb,
		EHighway:            v.EHighway,
		FuelCost:            v.F1FuelCost,
		FuelType:            v.F1FuelType,
		GhgScore:            v.F1GhgScore,
		MpgCity:             v.F1MpgCity,
		MpgCityUnadj:        v.F1MpgCityUnadj,
		MpgCityUnrounded:    v.F1MpgCityUnrounded,
		MpgComb:             v.F1MpgComb,
		MpgCombUnrounded:    v.F1MpgCombUnrounded,
		MpgHighway:          v.F1MpgHighway,
		MpgHighwayUnadj:     v.F1MpgHighwayUnadj,
		MpgHighwayUnrounded: v.F1MpgHighwayUnrounded,
		Range:               v.F1Range,
		PhevCDCity:          v.PhevCDCity,
		PhevCDComb:          v.PhevCDComb,
		PhevCDHighway:       v.PhevCDHighway,
		PhevMpgCity:         v.PhevMpgCity,
		PhevMpgComb:         v.PhevMpgComb,
		PhevMpgHighway:      v.PhevMpgHighway,
		PhevUFCity:          v.PhevUFCity,
		PhevUFComb:          v.PhevUFComb,
		PhevUFHighway:       v.PhevUFHighway,
	}
	return fuel
}

func buildFuel2(v *Vehicle) (fuel Fuel) {
	fuel = Fuel{
		BarrelsPerYear:      v.F2BarrelsPerYear,
		Co2:                 v.F2Co2,
		Co2Tailpipe:         v.F2Co2Tailpipe,
		FuelCost:            v.F2FuelCost,
		FuelType:            v.F2FuelType,
		GhgScore:            v.F2GhgScore,
		MpgCity:             v.F2MpgCity,
		MpgCityUnadj:        v.F2MpgCityUnadj,
		MpgCityUnrounded:    v.F2MpgCityUnrounded,
		MpgComb:             v.F2MpgComb,
		MpgCombUnrounded:    v.F2MpgCombUnrounded,
		MpgHighway:          v.F2MpgHighway,
		MpgHighwayUnadj:     v.F2MpgHighwayUnadj,
		MpgHighwayUnrounded: v.F2MpgHighwayUnrounded,
		Range:               v.F2Range,
		RangeCity:           v.F2RangeCity,
		RangeHighway:        v.F2RangeHighway,
	}
	return fuel
}

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
	MpgCityUnadj        float64 `json:"mpgCityUnadj,omitempty"`        // unadjusted city MPG for fuelType1; see the description of the EPA test procedures
	MpgCityUnrounded    float64 `json:"mpgCityUnrounded,omitempty"`    // unrounded city MPG for fuelType1 (2), (3)
	MpgComb             float64 `json:"-"`                             // combined MPG for fuelType1 (2)
	MpgCombUnrounded    float64 `json:"-"`                             // unrounded combined MPG for fuelType1 (2), (3)
	MpgHighway          float64 `json:"mpgHighway,omitempty"`          // highway MPG for fuelType1 (2)
	MpgHighwayUnadj     float64 `json:"mpgHighwayUnadj,omitempty"`     // unadjusted highway MPG for fuelType1; see the description of the EPA test procedures
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
