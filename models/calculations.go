package models

import "math"

var (
	CityShareDefault    int = 55
	HighwayShareDefault int = 45
	MilesPerYearDefault int = 15000
)

type DrivingProfile struct {
	CityShare    int `json:"cityShare"`
	HighwayShare int `json:"highwayShare"`
	MilesPerYear int `json:"milesPerYear"`
}

func calculateBarrelsPerYear(barrelsPer15000 float64, milesPerYear int) float64 {
	return toFixed(barrelsPer15000/15000.0*float64(milesPerYear), 2)
}

func calculateFuelCost(dollarsPerGallon float64, milesPerYear int, mpg float64) int {
	return int(dollarsPerGallon / mpg * float64(milesPerYear))
}

func calculateFuelCostElectricity(dollarsPerKwh float64, milesPerYear int, kwhPer100Miles float64) int {
	return int(float64(milesPerYear) / 100.0 * kwhPer100Miles * dollarsPerKwh)
}

func calculateEComb(cityKwh float64, cityShare int, highwayKwh float64, highwayShare int) float64 {
	return toFixed(float64(cityShare)/100.0*cityKwh+float64(highwayShare)/100.0*highwayKwh, 2)
}

func calculateMpgComb(cityMpg float64, cityShare int, highwayMpg float64, highwayShare int) float64 {
	return toFixed(float64(cityShare)/100.0*cityMpg+float64(highwayShare)/100.0*highwayMpg, 2)
}

func calculateRange(cityRange float64, cityShare int, highwayRange float64, highwayShare int) float64 {
	return toFixed(float64(cityShare)/100.0*cityRange+float64(highwayShare)/100.0*highwayRange, 2)
}

func fuelPriceFromName(name string, fp FuelPrices) (fuelPrice float64) {
	switch name {
	default:
		fuelPrice = 0.0
	case "Diesel":
		fuelPrice = fp.Diesel
	case "E85":
		fuelPrice = fp.E85
	case "Electricity":
		fuelPrice = fp.Electricity
	case "Midgrade Gasoline":
		fuelPrice = fp.GasMidgrade
	case "Natural Gas":
		fuelPrice = fp.CompressedNatGas
	case "Premium Gasoline":
		fuelPrice = fp.GasPremium
	case "Regular Gasoline":
		fuelPrice = fp.GasRegular
	}

	return fuelPrice
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
