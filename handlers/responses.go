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
	Vehicle models.Vehicle `json:"vehicle"`
}

type VehiclesResponse struct {
	Vehicles []models.Vehicle `json:"vehicles"`
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
