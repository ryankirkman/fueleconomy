package models

import "time"

type FuelPrices struct {
	ID               int       `xml:"-" db:"id, primaryKey"`      // Our Id
	Updated          time.Time `xml:"-" db:"updated, autoSet"`    // Our updated
	CompressedNatGas float64   `xml:"cng" db:"cng"`               // $ per gallon of gasoline equivalent (GGE) of compressed natural gas
	Diesel           float64   `xml:"diesel" db:"diesel"`         // $ per gallon of diesel
	E85              float64   `xml:"e85" db:"e85"`               // $ per gallon of E85
	Electricity      float64   `xml:"electric" db:"electricity"`  // $ per kw-hr of electricity
	GasMidgrade      float64   `xml:"midgrade" db:"gas_midgrade"` // $ per gallon of midgrade gasoline
	GasPremium       float64   `xml:"premium" db:"gas_premium"`   // $ per gallon of premium gasoline
	GasRegular       float64   `xml:"regular" db:"gas_regular"`   // $ per gallon of regular gasoline
	LiquidPropane    float64   `xml:"lpg" db:"liquid_propane"`    // $ per gallon of propane
}
