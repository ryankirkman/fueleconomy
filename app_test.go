package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/teasherm/fueleconomy/global"
	"github.com/teasherm/fueleconomy/handlers"
	"github.com/teasherm/fueleconomy/models"
	"github.com/teasherm/fueleconomy/workers"
)

// Utils
const GOPATH string = "/opt/go"
const ROOTPATH string = "/src/github.com/teasherm/fueleconomy"
const SQLITE_DB string = "/tmp/test.db"

func getPackagePath(packagePath string) string {
	return filepath.Join(GOPATH, ROOTPATH, packagePath)
}

func readFixture(fname string) (b []byte, err error) {
	filePath := filepath.Join(getPackagePath("models/fixtures"), fname)

	return ioutil.ReadFile(filePath)
}

// Mocks
var testServer *httptest.Server

type testFuelPricesFetcher struct{}

func (t testFuelPricesFetcher) Fetch(ignored string) ([]byte, error) {
	return readFixture("fuel_prices.xml")
}

type testVehiclesFetcher struct{}

func (t testVehiclesFetcher) Fetch(target string) ([]byte, error) {
	switch target {
	case "vehicles":
		return readFixture("vehicles.xml")
	case "emissions":
		return readFixture("emissions.xml")
	default:
		return make([]byte, 0), nil
	}
}

type flatVehicleResponse struct {
	Vehicle models.Vehicle `json:"vehicle"`
}

type flatVehiclesResponse struct {
	Vehicles []models.Vehicle `json:"vehicles"`
}

type DevNull struct{}

func (DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

// Tests
func TestVehicleGetOne(t *testing.T) {
	var vehicleGetOneUrl = fmt.Sprintf("%s/vehicle/1", testServer.URL)
	req, err := http.NewRequest("GET", vehicleGetOneUrl, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Vehicle get one http request not successful")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var fv flatVehicleResponse
	err = json.Unmarshal(body, &fv)
	if err != nil {
		t.Error(err.Error())
	}

	vehicle := fv.Vehicle

	if vehicle.EpaID != 1 {
		t.Error("Vehicle get one parsed epaID incorrectly")
	}

	if len(vehicle.EmissionsInfo) != 2 {
		t.Error("Vehicle get one parsed too few EmissionsInfoes")
	}

	if vehicle.EmissionsInfo[0].SalesArea != 3 {
		t.Error("Vehicle get one parsed EmissionsInfo salesArea incorrectly")
	}

	if len(vehicle.Fuels) != 1 {
		t.Error("Vehicle get one parsed too few fuels")
	}

	if vehicle.Fuels[0].MpgCity != 19.0 {
		t.Error("Vehicle get one parsed MPG City incorrectly")
	}
}

func TestVehicleGetManyExact(t *testing.T) {
	var vehicleGetManyUrl = fmt.Sprintf("%s/vehicles?year=1985", testServer.URL)
	req, err := http.NewRequest("GET", vehicleGetManyUrl, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Vehicles get many exact http request not successful")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var fvs flatVehiclesResponse
	err = json.Unmarshal(body, &fvs)
	if err != nil {
		t.Error(err.Error())
	}

	if len(fvs.Vehicles) != 1 {
		t.Error("Vehicle get many exact returned too few vehicles")
	}

	vehicle := fvs.Vehicles[0]

	if vehicle.EpaID != 1 {
		t.Error("Vehicles get many parsed epaID incorrectly")
	}

	if len(vehicle.EmissionsInfo) != 2 {
		t.Error("Vehicles get many parsed too few EmissionsInfoes")
	}

	if vehicle.EmissionsInfo[0].SalesArea != 3 {
		t.Error("Vehicles get many")
	}

	if len(vehicle.Fuels) != 1 {
		t.Error("Vehicles get many parsed too few fuels")
	}

	if vehicle.Fuels[0].MpgCity != 19.0 {
		t.Error("Vehicles get many parsed MPG City incorrectly")
	}
}

func TestVehicleGetManyFuzzy(t *testing.T) {
	var vehicleGetManyUrl = fmt.Sprintf("%s/vehicles?make=alfa", testServer.URL)
	req, err := http.NewRequest("GET", vehicleGetManyUrl, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		t.Error("Vehicles get many fuzzy not a 200")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var fvs flatVehiclesResponse
	err = json.Unmarshal(body, &fvs)
	if err != nil {
		t.Error(err.Error())
	}

	if len(fvs.Vehicles) != 1 {
		t.Error("Vehicle get many fuzzy failure")
	}

	vehicle := fvs.Vehicles[0]

	if vehicle.Make != "Alfa Romeo" {
		t.Error("EpaID parse failure")
	}
}

// Setup
func setup() (err error) {
	workRequest := workers.WorkRequest{
		Target:  "fuelprices",
		Fetcher: testFuelPricesFetcher{},
		Action:  workers.IngestFuelPrices}
	err = workRequest.DoWork()
	if err != nil {
		return err
	}

	workRequest = workers.WorkRequest{
		Target:  "vehicles",
		Fetcher: testVehiclesFetcher{},
		Action:  workers.IngestVehicles}
	err = workRequest.DoWork()
	if err != nil {
		return err
	}

	return nil
}

// Main
func TestMain(m *testing.M) {
	cmd := exec.Command("sqlite3", SQLITE_DB, "")
	_, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}

	cmd = exec.Command("sql-migrate", "up", "-env=test",
		fmt.Sprintf("-config=%s", getPackagePath("dbconfig.yml")))
	_, err = cmd.Output()
	if err != nil {
		os.Exit(1)
	}

	global.InitLogger(new(DevNull))
	global.InitDb("sqlite3", SQLITE_DB)
	testServer = httptest.NewServer(handlers.NewRouter())

	err = setup()
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()

	global.Db.Conn.Close()
	os.Remove(SQLITE_DB)
	os.Exit(code)
}
