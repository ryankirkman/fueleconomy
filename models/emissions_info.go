package models

import (
	"reflect"
	"time"
)

func NewEmissionsInfoFromRaw(r *RawEmissionsInfo) (*EmissionsInfo, error) {
	e := EmissionsInfo{}
	ev := reflect.ValueOf(&e)
	rv := reflect.ValueOf(r).Elem()
	for i := 0; i < rv.NumField(); i++ {
		rStructFld := rv.Type().Field(i)
		rFldValue := rv.Field(i)
		err := parseStringToOutStruct(rFldValue, rStructFld, ev)
		if err != nil {
			return nil, err
		}
	}
	return &e, nil
}

type EmissionsInfo struct {
	ID              int       `db:"id, primaryKey" json:"-"`                            // Our ID
	Updated         time.Time `db:"updated, autoSet" json:"updated,omitempty"`          // Our update timestamp
	EmissionStdCode string    `db:"emission_std_code" json:"emissionStdCode,omitempty"` // Vehicle Emission Standard Code
	EmissionStdTxt  string    `db:"emission_std_txt" json:"emissionsStdTxt,omitempty"`  // Vehicle Emission Standard
	EngineFamilyID  string    `db:"engine_family_id" json:"engineFamilyId,omitempty"`   // engine family ID
	EpaID           int       `db:"epa_id" json:"-"`                                    // vehicle record ID (links emission data to the vehicle record)
	F1SmogRating    float64   `db:"f1_smog_rating" json:"f1SmogRating,omitempty"`       // EPA 1/10 smog rating for fuelType1
	F2SmogRating    float64   `db:"f2_smog_rating" json:"f2SmogRating,omitempty"`       // EPA 1/10 smog rating for fuelType2
	SalesArea       int       `db:"sales_area" json:"salesArea,omitempty"`              // EPA sales area code
	SmartwayScore   int       `db:"smartway_score" json:"smartwayScore,omitempty"`      // SmartWay Code
}

type RawEmissionsInfoOuter struct {
	RawEmissionsInfoes []RawEmissionsInfo `xml:"emissionsInfo"`
}

type RawEmissionsInfo struct {
	EmissionStdCode string `xml:"standard"`      // Vehicle Emission Standard Code
	EmissionStdTxt  string `xml:"stdText"`       // Vehicle Emission Standard
	EngineFamilyID  string `xml:"efid"`          // engine family ID
	EpaID           string `xml:"id"`            // vehicle record ID (links emission data to the vehicle record)
	F1SmogRating    string `xml:"score"`         // EPA 1/10 smog rating for fuelType1
	F2SmogRating    string `xml:"scoreAlt"`      // EPA 1/10 smog rating for fuelType2
	SalesArea       string `xml:"salesArea"`     // EPA sales area code
	SmartwayScore   string `xml:"smartwayScore"` // SmartWay Code
}
