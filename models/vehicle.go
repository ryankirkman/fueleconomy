package models

import (
	"reflect"
	"time"
)

func NewVehicleFromRaw(raw *RawVehicle) (*Vehicle, error) {
	vehicle := Vehicle{}
	vehicleVal := reflect.ValueOf(&vehicle)
	rawVal := reflect.ValueOf(raw).Elem()
	for i := 0; i < rawVal.NumField(); i++ {
		rawStructFld := rawVal.Type().Field(i)
		rawFldValue := rawVal.Field(i)
		err := parseStringToOutStruct(rawFldValue, rawStructFld, vehicleVal)
		if err != nil {
			return nil, err
		}
	}
	return &vehicle, nil
}

type Vehicle struct {
	ID                    int             `db:"id, primaryKey" json:"-"`                                      // Our ID
	Updated               time.Time       `db:"updated, autoSet" json:"updated"`                              // Our updated timestamp
	AtvType               string          `db:"atv_type" json:"atvType,omitempty"`                            // type of alternative fuel or advanced technology vehicle
	ChargeTime120V        float64         `db:"charge_time_120v" json:"chargeTime120V,omitempty"`             // time to charge an electric vehicle in hours at 120 V
	ChargeTime240V        float64         `db:"charge_time_240v" json:"chargeTime240V,omitempty"`             // time to charge an electric vehicle in hours at 240 V
	ChargeTime240Vb       float64         `db:"charge_time_240vb" json:"chargeTime240Vb,omitempty"`           // time to charge an electric vehicle in hours at 240 V using the alternate charger
	Charger240VDscr       string          `db:"charger_240v_dscr" json:"charger240VDscr,omitempty"`           // electric vehicle charger description
	Charger240VbDscr      string          `db:"charger_240vb_dscr" json:"charger240VbDscr,omitempty"`         // electric vehicle alternate charger description
	Cylinders             int             `db:"cylinders" json:"cylinders,omitempty"`                         // engine cylinders
	DriveAxleType         string          `db:"drive_axle_type" json:"driveAxleType,omitempty"`               // drive axle type
	ECity                 float64         `db:"e_city" json:"-"`                                              // city electricity consumption in kw-hrs/100 miles
	EComb                 float64         `db:"e_comb" json:"-"`                                              // combined electricity consumption in kw-hrs/100 miles
	EHighway              float64         `db:"e_highway" json:"-"`                                           // highway electricity consumption in kw-hrs/100 miles
	EMotor                string          `db:"e_motor" json:"eMotor,omitempty"`                              // electric motor (kw-hrs)
	EmissionsInfo         []EmissionsInfo `db:"-" json:"emissionsInfo,omitempty"`                             // Emissions Info from separate table, joined at application level
	EngDisplacement       float64         `db:"eng_displacement" json:"engDisplacement,omitempty"`            // engine displacement in liters
	EngDscr               string          `db:"eng_dscr" json:"engDscr,omitempty"`                            // engine descriptor; see http://www.fueleconomy.gov/feg/findacarhelp.shtml#engine
	EngID                 int             `db:"eng_id" json:"engID,omitempty"`                                // EPA model type index
	EpaCreatedOn          time.Time       `db:"epa_created_on" json:"epaCreatedOn,omitempty"`                 // date the vehicle record was created (ISO 8601 format)
	EpaID                 int             `db:"epa_id" json:"epaID,omitempty"`                                // vehicle record id
	EpaModifiedOn         time.Time       `db:"epa_modified_on" json:"epaModifiedOn,omitempty"`               // date the vehicle record was last modified (ISO 8601 format)
	F1BarrelsPerYear      float64         `db:"f1_barrels_per_year" json:"-"`                                 // annual petroleum consumption in barrels for fuelType1 (1)
	F1Co2                 float64         `db:"f1_co2" json:"-"`                                              // tailpipe CO2 in grams/mile for fuelType1 (5)
	F1Co2Tailpipe         float64         `db:"f1_co2_tailpipe" json:"-"`                                     // tailpipe CO2 in grams/mile for fuelType1 (5)
	F1FuelCost            int             `db:"f1_fuel_cost" json:"-"`                                        // annual fuel cost for fuelType1 ($) (7)
	F1FuelType            string          `db:"f1_fuel_type" json:"-"`                                        // fuel type 1. For single fuel vehicles, this will be the only fuel. For dual fuel vehicles, this will be the conventional fuel.
	F1GhgScore            int             `db:"f1_ghg_score" json:"-"`                                        // EPA GHG score (-1 = Not available)
	F1MpgCity             float64         `db:"f1_mpg_city" json:"-"`                                         // city MPG for fuelType1 (2)
	F1MpgCityUnadj        float64         `db:"f1_mpg_city_unadj" json:"-"`                                   // unadjusted city MPG for fuelType1; see the description of the EPA test procedures
	F1MpgCityUnrounded    float64         `db:"f1_mpg_city_unrounded" json:"-"`                               // unrounded city MPG for fuelType1 (2), (3)
	F1MpgComb             float64         `db:"f1_mpg_comb" json:"-"`                                         // combined MPG for fuelType1 (2)
	F1MpgCombUnrounded    float64         `db:"f1_mpg_comb_unrounded" json:"-"`                               // unrounded combined MPG for fuelType1 (2), (3)
	F1MpgHighway          float64         `db:"f1_mpg_highway" json:"-"`                                      // highway MPG for fuelType1 (2)
	F1MpgHighwayUnadj     float64         `db:"f1_mpg_highway_unadj" json:"-"`                                // unadjusted highway MPG for fuelType1; see the description of the EPA test procedures
	F1MpgHighwayUnrounded float64         `db:"f1_mpg_highway_unrounded" json:"-"`                            // unrounded highway MPG for fuelType1 (2), (3)
	F1Range               float64         `db:"f1_range" json:"-"`                                            // EPA range for fuelType2
	F2BarrelsPerYear      float64         `db:"f2_barrels_per_year" json:"-"`                                 // annual petroleum consumption in barrels for fuelType2 (1)
	F2Co2                 float64         `db:"f2_co2" json:"-"`                                              // tailpipe CO2 in grams/mile for fuelType2 (5)
	F2Co2Tailpipe         float64         `db:"f2_co2_tailpipe" json:"-"`                                     // tailpipe CO2 in grams/mile for fuelType2 (5)
	F2FuelCost            int             `db:"f2_fuel_cost" json:"-"`                                        // annual fuel cost for fuelType2 ($) (7)
	F2FuelType            string          `db:"f2_fuel_type" json:"-"`                                        // fuel type 2. For dual fuel vehicles, this will be the alternative fuel (e.g. E85, Electricity, CNG, LPG). For single fuel vehicles, this field is not used
	F2GhgScore            int             `db:"f2_ghg_score" json:"-"`                                        // EPA GHG score for dual fuel vehicle running on the alternative fuel (//1 = Not available)
	F2MpgCity             float64         `db:"f2_mpg_city" json:"-"`                                         // city MPG for fuelType2 (2)
	F2MpgCityUnadj        float64         `db:"f2_mpg_city_unadj" json:"-"`                                   // unadjusted city MPG for fuelType2; see the description of the EPA test procedures
	F2MpgCityUnrounded    float64         `db:"f2_mpg_city_unrounded" json:"-"`                               // unrounded city MPG for fuelType2 (2), (3)
	F2MpgComb             float64         `db:"f2_mpg_comb" json:"-"`                                         // combined MPG for fuelType2 (2)
	F2MpgCombUnrounded    float64         `db:"f2_mpg_comb_unrounded" json:"-"`                               // unrounded combined MPG for fuelType2 (2), (3)
	F2MpgHighway          float64         `db:"f2_mpg_highway" json:"-"`                                      // highway MPG for fuelType2 (2)
	F2MpgHighwayUnrounded float64         `db:"f2_mpg_highway_unrounded" json:"-"`                            // unrounded highway MPG for fuelType2 (2),(3)
	F2MpgHighwayUnadj     float64         `db:"f2_mpg_highway_unadj" json:"-"`                                // unadjusted highway MPG for fuelType2; see the description of the EPA test procedures
	F2Range               float64         `db:"f2_range" json:"-"`                                            // EPA range for fuelType2
	F2RangeCity           float64         `db:"f2_range_city" json:"-"`                                       // EPA city range for fuelType2
	F2RangeHighway        float64         `db:"f2_range_highway" json:"-"`                                    // EPA highway range for fuelType2
	FuelEconomyScore      float64         `db:"fuel_economy_score" json:"fuelEconomyScore,omitempty"`         // EPA Fuel Economy Score (-1 = Not available)
	Fuels                 []Fuel          `db:"-" json:"fuels"`                                               // Processed slice of fuels (min length 1, max length 2)
	FuelType              string          `db:"fuel_type" json:"fuelType,omitempty"`                          // fuel type with fuelType1 and fuelType2 (if applicable)
	HasMpgData            bool            `db:"has_mpg_data" json:"-"`                                        // Parsed boolean
	HasStartStop          bool            `db:"start_stop" json:"hasStartStop,omitempty"`                     // Parsed boolean
	HasSupercharger       bool            `db:"has_supercharger" json:"hasSupercharger,omitempty"`            // Parsed boolean
	HasTurbocharger       bool            `db:"has_turbocharger" json:"hasTurbocharger,omitempty"`            // Parsed boolean
	IsGuzzler             bool            `db:"is_guzzler" json:"isGuzzler,omitempty"`                        // Parsed boolean
	IsPhevBlended         bool            `db:"is_phev_blended" json:"isPhevBlended,omitempty"`               // Parsed boolean
	LuggageVolume2Door    int             `db:"luggage_volume_2door" json:"luggageVolume2Door,omitempty"`     // 2 door luggage volume (cubic feet) (8)
	LuggageVolume4Door    int             `db:"luggage_volume_4door" json:"luggageVolume4Door,omitempty"`     // 4 door luggage volume (cubic feet) (8)
	LuggageVolumeHatch    int             `db:"luggage_volume_hatch" json:"luggageVolumeHatch,omitempty"`     // hatchback luggage volume (cubic feet) (8)
	Make                  string          `db:"make" json:"make,omitempty"`                                   // manufacturer (division)
	ManufacturerCode      string          `db:"manufacturer_code" json:"manufacturerCode,omitempty"`          // 3-character manufacturer code
	Model                 string          `db:"model" json:"model,omitempty"`                                 // model name (carline)
	PassengerVolume2Door  int             `db:"passenger_volume_2door" json:"passengerVolume2Door,omitempty"` // 2-door passenger volume (cubic feet) (8)
	PassengerVolume4Door  int             `db:"passenger_volume_4door" json:"passengerVolume4Door,omitempty"` // 4-door passenger volume (cubic feet) (8)
	PassengerVolumeHatch  int             `db:"passenger_volume_hatch" json:"passengerVolumeHatch,omitempty"` // hatchback passenger volume (cubic feet) (8)
	PhevCDCity            float64         `db:"phev_cd_city" json:"-"`                                        // city gasoline consumption (gallons/100 miles) in charge depleting mode (4)
	PhevCDComb            float64         `db:"phev_cd_comb" json:"-"`                                        // combined gasoline consumption (gallons/100 miles) in charge depleting mode (4)
	PhevCDHighway         float64         `db:"phev_cd_highway" json:"-"`                                     // highway gasoline consumption (gallons/100miles) in charge depleting mode (4)
	PhevMpgCity           float64         `db:"phev_mpg_city" json:"-"`                                       // EPA composite gasoline-electricity city MPGe for plug-in hybrid vehicles
	PhevMpgComb           float64         `db:"phev_mpg_comb" json:"-"`                                       // EPA composite gasoline-electricity combined city-highway MPGe for plug-in hybrid vehicles
	PhevMpgHighway        float64         `db:"phev_mpg_highway" json:"-"`                                    // EPA composite gasoline-electricity highway MPGe for plug-in hybrid vehicles
	PhevUFCity            float64         `db:"phev_uf_city" json:"-"`                                        // EPA city utility factor (share of electricity) for PHEV
	PhevUFComb            float64         `db:"phev_uf_comb" json:"-"`                                        // EPA combined utility factor (share of electricity) for PHEV
	PhevUFHighway         float64         `db:"phev_uf_highway" json:"-"`                                     // EPA highway utility factor (share of electricity) for PHEV
	SizeClass             string          `db:"size_class" json:"sizeClass,omitempty"`                        // EPA vehicle size class
	TransDscr             string          `db:"trans_dscr" json:"transDscr,omitempty"`                        // transmission descriptor; see http://www.fueleconomy.gov/feg/findacarhelp.shtml#trany
	Transition            string          `db:"transition" json:"transition,omitempty"`                       // transmission
	Year                  int             `db:"year" json:"year,omitempty"`                                   // model year
	YouSaveSpend          int             `db:"you_save_spend" json:"-"`                                      // you save/spend over 5 years compared to an average car ($). Savings are positive; a greater amount spent yields a negative number. For dual fuel vehicles, this is the cost savings for gasoline
}

type RawVehiclesOuter struct {
	RawVehicles []RawVehicle `xml:"vehicle"`
}

type RawVehicle struct {
	AtvType               string `xml:"atvtype"`                                    // type of alternative fuel or advanced technology vehicle
	ChargeTime120V        string `xml:"charge120"`                                  // time to charge an electric vehicle in hours at 120 V
	ChargeTime240V        string `xml:"charge240"`                                  // time to charge an electric vehicle in hours at 240 V
	ChargeTime240Vb       string `xml:"charge240b"`                                 // time to charge an electric vehicle in hours at 240 V using the alternate charger
	Charger240VDscr       string `xml:"c240Dscr"`                                   // electric vehicle charger description
	Charger240VbDscr      string `xml:"c240bDscr"`                                  // electric vehicle alternate charger description
	Cylinders             string `xml:"cylinders"`                                  // engine cylinders
	DriveAxleType         string `xml:"drive"`                                      // drive axle type
	ECity                 string `xml:"cityE"`                                      // city electricity consumption in kw-hrs/100 miles
	EComb                 string `xml:"combE"`                                      // combined electricity consumption in kw-hrs/100 miles
	EHighway              string `xml:"highwayE"`                                   // highway electricity consumption in kw-hrs/100 miles
	EMotor                string `xml:"evMotor"`                                    // electric motor (kw-hrs)
	EngDisplacement       string `xml:"displ"`                                      // engine displacement in liters
	EngDscr               string `xml:"eng_dscr"`                                   // engine descriptor; see http://www.fueleconomy.gov/feg/findacarhelp.shtml#engine
	EngID                 string `xml:"engId"`                                      // EPA model type index
	EpaCreatedOn          string `xml:"createdOn"`                                  // date the vehicle record was created (ISO 8601 format)
	EpaID                 string `xml:"id"`                                         // vehicle record id
	EpaModifiedOn         string `xml:"modifiedOn"`                                 // date the vehicle record was last modified (ISO 8601 format)
	F1BarrelsPerYear      string `xml:"barrels08"`                                  // annual petroleum consumption in barrels for fuelType1 (1)
	F1Co2                 string `xml:"co2"`                                        // tailpipe CO2 in grams/mile for fuelType1 (5)
	F1Co2Tailpipe         string `xml:"co2TailpipeGpm"`                             // tailpipe CO2 in grams/mile for fuelType1 (5)
	F1FuelCost            string `xml:"fuelCost08"`                                 // annual fuel cost for fuelType1 ($) (7)
	F1FuelType            string `xml:"fuelType1"`                                  // fuel type 1. For single fuel vehicles, this will be the only fuel. For dual fuel vehicles, this will be the conventional fuel.
	F1GhgScore            string `xml:"ghgScore"`                                   // EPA GHG score (-1 = Not available)
	F1MpgCity             string `xml:"city08"`                                     // city MPG for fuelType1 (2)
	F1MpgCityUnadj        string `xml:"UCity"`                                      // unadjusted city MPG for fuelType1; see the description of the EPA test procedures
	F1MpgCityUnrounded    string `xml:"city08U"`                                    // unrounded city MPG for fuelType1 (2), (3)
	F1MpgComb             string `xml:"comb08"`                                     // combined MPG for fuelType1 (2)
	F1MpgCombUnrounded    string `xml:"comb08U"`                                    // unrounded combined MPG for fuelType1 (2), (3)
	F1MpgHighway          string `xml:"highway08"`                                  // highway MPG for fuelType1 (2)
	F1MpgHighwayUnadj     string `xml:"UHighway"`                                   // unadjusted highway MPG for fuelType1; see the description of the EPA test procedures
	F1MpgHighwayUnrounded string `xml:"highway08U"`                                 // unrounded highway MPG for fuelType1 (2), (3)
	F1Range               string `xml:"range"`                                      // EPA range for fuelType2
	F2BarrelsPerYear      string `xml:"barrelsA08"`                                 // annual petroleum consumption in barrels for fuelType2 (1)
	F2Co2                 string `xml:"co2A"`                                       // tailpipe CO2 in grams/mile for fuelType2 (5)
	F2Co2Tailpipe         string `xml:"co2TailpipeAGpm"`                            // tailpipe CO2 in grams/mile for fuelType2 (5)
	F2FuelCost            string `xml:"fuelCostA08"`                                // annual fuel cost for fuelType2 ($) (7)
	F2FuelType            string `xml:"fuelType2"`                                  // fuel type 2. For dual fuel vehicles, this will be the alternative fuel (e.g. E85, Electricity, CNG, LPG). For single fuel vehicles, this field is not used
	F2GhgScore            string `xml:"ghgScoreA"`                                  // EPA GHG score for dual fuel vehicle running on the alternative fuel (//1 = Not available)
	F2MpgCity             string `xml:"cityA08"`                                    // city MPG for fuelType2 (2)
	F2MpgCityUnadj        string `xml:"UCityA"`                                     // unadjusted city MPG for fuelType2; see the description of the EPA test procedures
	F2MpgCityUnrounded    string `xml:"cityA08U"`                                   // unrounded city MPG for fuelType2 (2), (3)
	F2MpgComb             string `xml:"combA08"`                                    // combined MPG for fuelType2 (2)
	F2MpgCombUnrounded    string `xml:"combA08U"`                                   // unrounded combined MPG for fuelType2 (2), (3)
	F2MpgHighway          string `xml:"highwayA08"`                                 // highway MPG for fuelType2 (2)
	F2MpgHighwayUnrounded string `xml:"highwayA08U"`                                // unrounded highway MPG for fuelType2 (2),(3)
	F2MpgHighwayUnadj     string `xml:"UHighwayA"`                                  // unadjusted highway MPG for fuelType2; see the description of the EPA test procedures
	F2Range               string `xml:"rangeA"`                                     // EPA range for fuelType2
	F2RangeCity           string `xml:"rangeCityA"`                                 // EPA city range for fuelType2
	F2RangeHighway        string `xml:"rangeHwyA"`                                  // EPA highway range for fuelType2
	FuelEconomyScore      string `xml:"feScore"`                                    // EPA Fuel Economy Score (-1 = Not available)
	FuelType              string `xml:"fuelType"`                                   // fuel type with fuelType1 and fuelType2 (if applicable)
	Guzzler               string `xml:"guzzler" parseBool:"IsGuzzler,T,G"`          // if G or T, this vehicle is subject to the gas guzzler tax
	LuggageVolume2Door    string `xml:"lv2"`                                        // 2 door luggage volume (cubic feet) (8)
	LuggageVolume4Door    string `xml:"lv4"`                                        // 4 door luggage volume (cubic feet) (8)
	LuggageVolumeHatch    string `xml:"hlv"`                                        // hatchback luggage volume (cubic feet) (8)
	Make                  string `xml:"make"`                                       // manufacturer (division)
	ManufacturerCode      string `xml:"mfrCode"`                                    // 3-character manufacturer code
	Model                 string `xml:"model"`                                      // model name (carline)
	MpgData               string `xml:"mpgData" parseBool:"HasMpgData,Y"`           // has My MPG data; see yourMpgVehicle and yourMpgDriverVehicle
	PassengerVolume2Door  string `xml:"pv2"`                                        // 2-door passenger volume (cubic feet) (8)
	PassengerVolume4Door  string `xml:"pv4"`                                        // 4-door passenger volume (cubic feet) (8)
	PassengerVolumeHatch  string `xml:"hpv"`                                        // hatchback passenger volume (cubic feet) (8)
	PhevBlended           string `xml:"phevBlended" parseBool:"IsPhevBlended,true"` // if true, this vehicle operates on a blend of gasoline and electricity in charge depleting mode
	PhevCDCity            string `xml:"cityCD"`                                     // city gasoline consumption (gallons/100 miles) in charge depleting mode (4)
	PhevCDComb            string `xml:"combinedCD"`                                 // combined gasoline consumption (gallons/100 miles) in charge depleting mode (4)
	PhevCDHighway         string `xml:"highwayCD"`                                  // highway gasoline consumption (gallons/100miles) in charge depleting mode (4)
	PhevMpgCity           string `xml:"phevCity"`                                   // EPA composite gasoline-electricity city MPGe for plug-in hybrid vehicles
	PhevMpgComb           string `xml:"phevComb"`                                   // EPA composite gasoline-electricity combined city-highway MPGe for plug-in hybrid vehicles
	PhevMpgHighway        string `xml:"phevHwy"`                                    // EPA composite gasoline-electricity highway MPGe for plug-in hybrid vehicles
	PhevUFCity            string `xml:"cityUF"`                                     // EPA city utility factor (share of electricity) for PHEV
	PhevUFComb            string `xml:"combinedUF"`                                 // EPA combined utility factor (share of electricity) for PHEV
	PhevUFHighway         string `xml:"highwayUF"`                                  // EPA highway utility factor (share of electricity) for PHEV
	SizeClass             string `xml:"VClass"`                                     // EPA vehicle size class
	StartStop             string `xml:"startStop" parseBool:"HasStartStop,Y"`       // vehicle has start-stop technology (Y, N, or blank for older vehicles)
	Supercharger          string `xml:"sCharger" parseBool:"HasSupercharger,S"`     // if S, this vehicle is supercharged
	TransDscr             string `xml:"trans_dscr"`                                 // transmission descriptor; see http://www.fueleconomy.gov/feg/findacarhelp.shtml#trany
	Transition            string `xml:"trany"`                                      // transmission
	Turbocharger          string `xml:"tCharger" parseBool:"HasTurbocharger,T"`     // if T, this vehicle is turbocharged
	Year                  string `xml:"year"`                                       // model year
	YouSaveSpend          string `xml:"youSaveSpend"`                               // you save/spend over 5 years compared to an average car ($). Savings are positive; a greater amount spent yields a negative number. For dual fuel vehicles, this is the cost savings for gasoline
}
