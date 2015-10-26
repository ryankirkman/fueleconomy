package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teasherm/fueleconomy/models"
)

type SimpleResponse struct {
	Message string `json:"message"`
}

type VehicleResponse struct {
	Profile models.DrivingProfile `json:"profile"`
	Vehicle models.Vehicle        `json:"vehicle"`
}

type VehiclesResponse struct {
	Meta     PageInfo              `json:"meta"`
	Profile  models.DrivingProfile `json:"profile"`
	Vehicles []models.Vehicle      `json:"vehicles"`
}

func sendErrorJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	js, _ := json.Marshal(SimpleResponse{message})
	w.WriteHeader(code)
	w.Write(js)
}

func sendJSON(w http.ResponseWriter, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
