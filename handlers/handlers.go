package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/models"
	"github.com/teasherm/fueleconomy/srm"
	"github.com/teasherm/fueleconomy/workers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/health_check", HealthCheck).Methods("GET")
	r.HandleFunc("/ingest/{target}", Ingest).Methods("GET")
	r.HandleFunc("/vehicle/{id:[0-9]+}", VehicleGetOne).Methods("GET")
	r.HandleFunc("/vehicles", VehicleGetMany).Methods("GET")

	return r
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := global.Db.Conn.Ping()
	checkErr(err, w)

	js, err := json.Marshal(SimpleResponse{"Healthy!"})
	checkErr(err, w)
	sendJSON(w, js)
}

func Index(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile("/opt/go/src/github.com/teasherm/fueleconomy/frontend/index.html")
	checkErr(err, w)
	fmt.Fprintf(w, string(body))
}

func Ingest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	target := vars["target"]
	work, err := workers.GenerateWorkRequest(target)
	checkErr(err, w)
	workers.WorkQueue <- work

	js, err := json.Marshal(SimpleResponse{fmt.Sprintf("Ingest kicked off for: %s", target)})
	checkErr(err, w)
	sendJSON(w, js)
}

func VehicleGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	queryVals := r.URL.Query()
	profile := getProfileFromQueryVals(queryVals)

	v := models.Vehicle{}
	eis := make([]models.EmissionsInfo, 0)
	query := fmt.Sprintf("SELECT * FROM vehicles WHERE epa_id = %s",
		global.Db.Dialect.Placeholder(1))
	err := global.Db.SelectOne(&v, query, id)
	checkErr(err, w)

	fp := getMostRecentFuelPrices()
	v.Fuels = models.CalculateFuelData(&v, profile, fp)

	query = fmt.Sprintf("SELECT * FROM emissions_info WHERE epa_id = %s",
		global.Db.Dialect.Placeholder(1))
	err = global.Db.SelectMany(&eis, query, id)
	checkErr(err, w)
	v.EmissionsInfo = eis

	js, err := json.Marshal(VehicleResponse{profile, v})
	checkErr(err, w)
	sendJSON(w, js)
}

var (
	// Querties correspond to exact matches
	// SQL: WHERE col = query
	ExactParams []searchParam = []searchParam{
		searchParam{name: "year", converter: intConverter},
	}
	// Queries are subject to case insensitive matching with wildcard tails
	// SQL: WHERE lower(col) LIKE 'lower(query)%'
	FuzzyParams []string = []string{"make", "model"}
)

func VehicleGetMany(w http.ResponseWriter, r *http.Request) {
	// Parse querystring parameters and make sql query builder
	queryVals := r.URL.Query()
	profile := getProfileFromQueryVals(queryVals)
	page := getPageFromQueryVals(queryVals, r.URL)
	queryBuilder := &srm.QueryBuilder{
		Db:         global.Db,
		Table:      "vehicles",
		Limit:      page.PageLength,
		Offset:     page.PageLength * (page.PageNo - 1),
		WhereExact: extractSearchParams(queryVals, ExactParams),
		WhereFuzzy: extractStringParams(queryVals, FuzzyParams),
	}

	// Get results count
	query, vals := queryBuilder.BuildCount()
	resultCount, err := global.Db.SelectInt(query, vals...)
	checkErr(err, w)
	page.Fill(queryVals, resultCount)

	// Query for page of vehicles
	query, vals = queryBuilder.BuildSelect()
	vs := make([]models.Vehicle, 0)
	err = global.Db.SelectMany(&vs, query, vals...)
	checkErr(err, w)

	// Calculate fuel data on vehicles
	fp := getMostRecentFuelPrices()
	epaIdsQuery, epaIds, epaIdToIdx := calculateFuelDataForAndCollectEpaIdsFromVehicles(
		&vs, profile, fp)

	// Query for emissions info and append to vehicles
	eis := make([]models.EmissionsInfo, 0)
	query = fmt.Sprintf("SELECT * FROM emissions_info WHERE epa_id IN (%s)", epaIdsQuery)
	global.Db.SelectMany(&eis, query, epaIds...)
	for _, ei := range eis {
		v := &vs[epaIdToIdx[ei.EpaID]]
		v.EmissionsInfo = append(v.EmissionsInfo, ei)
	}

	// Send response
	js, err := json.Marshal(VehiclesResponse{*page, profile, vs})
	checkErr(err, w)
	sendJSON(w, js)
}

func calculateFuelDataForAndCollectEpaIdsFromVehicles(vehicles *[]models.Vehicle,
	profile models.DrivingProfile, fp models.FuelPrices) (epaIdsQuery string, epaIds []interface{},
	epaIdToIdx map[int]int) {

	queryBuff := bytes.Buffer{}
	epaIdToIdx = make(map[int]int, 0)
	vs := *vehicles
	for i := 0; i < len(vs); i++ {
		v := &vs[i]
		epaIdToIdx[v.EpaID] = i
		epaIds = append(epaIds, v.EpaID)
		if i > 0 {
			queryBuff.WriteString(", ")
		}
		queryBuff.WriteString(global.Db.Dialect.Placeholder(i + 1))
		v.Fuels = models.CalculateFuelData(v, profile, fp)
	}
	return queryBuff.String(), epaIds, epaIdToIdx
}
