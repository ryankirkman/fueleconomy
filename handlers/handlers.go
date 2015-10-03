package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/models"
	"github.com/teasherm/fueleconomy/srm"
	"github.com/teasherm/fueleconomy/workers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health_check", HealthCheck).Methods("GET")
	r.HandleFunc("/ingest/{target}", Ingest).Methods("GET")
	r.HandleFunc("/vehicle/{id:[0-9]+}", VehicleGetOne).Methods("GET")
	r.HandleFunc("/vehicles", VehicleGetMany).Methods("GET")

	return r
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := global.Db.Conn.Ping()
	if err != nil {
		sendErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, _ := json.Marshal(SimpleResponse{"Healthy!"})
	sendJSON(w, js)
}

func Ingest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	target := vars["target"]
	work, err := workers.GenerateWorkRequest(target)
	if err != nil {
		sendErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	workers.WorkQueue <- work
	js, err := json.Marshal(SimpleResponse{fmt.Sprintf("Ingest kicked off for: %s", target)})
	if err != nil {
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
	sendJSON(w, js)
}

func VehicleGetOne(w http.ResponseWriter, r *http.Request) {
	v := models.Vehicle{}
	eis := make([]models.EmissionsInfo, 0)
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	query := fmt.Sprintf("SELECT * FROM vehicles WHERE epa_id = %s", global.Db.Dialect.Placeholder(1))
	err := global.Db.SelectOne(&v, query, id)
	if err != nil {
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
	query = fmt.Sprintf("SELECT * FROM emissions_info WHERE epa_id = %s", global.Db.Dialect.Placeholder(1))
	err = global.Db.SelectMany(&eis, query, id)
	if err != nil {
		global.Logger.Println("Error: ", err)
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
	v.EmissionsInfo = eis
	v.Fuels = models.ConstructFuelData(&v)
	js, err := json.Marshal(VehicleResponse{v})
	if err != nil {
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
	sendJSON(w, js)
}

// Multiple get paramters
var (
	// Querties correspond to exact matches
	// SQL: WHERE col = query
	MATCH_EXACT []string = []string{"year"}
	// Queries are subject to case insensitive matching with wildcard tails
	// SQL: WHERE lower(col) LIKE 'lower(query)%'
	MATCH_FUZZY []string = []string{"make", "model"}
	// Default page length
	PAGE_SIZE_DEFAULT int = 10
	PAGE_SIZE_MAX     int = 100
)

func extractParamsFromQueryString(queryString url.Values, params []string) (res map[string]string) {
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

func VehicleGetMany(w http.ResponseWriter, r *http.Request) {
	queryString := r.URL.Query()
	queryBuilder := srm.QueryBuilder{
		Db:    global.Db,
		Table: "vehicles",
	}
	whereExact := extractParamsFromQueryString(queryString, MATCH_EXACT)
	whereFuzzy := extractParamsFromQueryString(queryString, MATCH_FUZZY)

	pageNo := 1
	pageQuery := queryString.Get("page")
	if parsedPageQuery, err := strconv.Atoi(pageQuery); err != nil {
		pageNo = parsedPageQuery
	}

	pageSize := PAGE_SIZE_DEFAULT
	pageSizeQuery := queryString.Get("pageSize")
	if parsedPageSizeQuery, err := strconv.Atoi(pageSizeQuery); err != nil {
		pageSize = maxInt(parsedPageSizeQuery, PAGE_SIZE_MAX)
	}

	query, vals := queryBuilder.BuildSelect(pageSize, pageSize*(pageNo-1),
		whereExact, whereFuzzy)

	vs := make([]models.Vehicle, 0)

	err := global.Db.SelectMany(&vs, query, vals...)
	if err != nil {
		global.Logger.Println(err)
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
	idToIdx := make(map[int]int)
	ids := make([]interface{}, 0)
	buff := bytes.Buffer{}
	for i := 0; i < len(vs); i++ {
		v := &vs[i]
		idToIdx[v.EpaID] = i
		ids = append(ids, v.EpaID)
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(global.Db.Dialect.Placeholder(i + 1))
		v.Fuels = models.ConstructFuelData(v)
	}
	query = fmt.Sprintf("SELECT * FROM emissions_info WHERE epa_id IN (%s)", buff.String())
	eis := make([]models.EmissionsInfo, 0)
	global.Db.SelectMany(&eis, query, ids...)
	for _, ei := range eis {
		v := &vs[idToIdx[ei.EpaID]]
		v.EmissionsInfo = append(v.EmissionsInfo, ei)
	}
	js, err := json.Marshal(VehiclesResponse{vs})
	if err != nil {
		global.Logger.Println(err)
		sendErrorJSON(w, "Server error", http.StatusInternalServerError)
		return
	}
	sendJSON(w, js)
}
