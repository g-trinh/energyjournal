package energy

import (
	"encoding/json"
	"net/http"
)

type EnergyLevelsResponse struct {
	Date      string `json:"date"`
	Physical  int    `json:"physical"`
	Mental    int    `json:"mental"`
	Emotional int    `json:"emotional"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
