package energy

import (
	"encoding/json"
	"net/http"
)

type EnergyLevelsResponse struct {
	Date               string `json:"date"`
	Physical           int    `json:"physical"`
	Mental             int    `json:"mental"`
	Emotional          int    `json:"emotional"`
	SleepQuality       int    `json:"sleepQuality,omitempty"`
	StressLevel        int    `json:"stressLevel,omitempty"`
	PhysicalActivity   string `json:"physicalActivity,omitempty" enums:"none,light,moderate,intense"`
	Nutrition          string `json:"nutrition,omitempty" enums:"poor,average,good,excellent"`
	SocialInteractions string `json:"socialInteractions,omitempty" enums:"negative,neutral,positive"`
	TimeOutdoors       string `json:"timeOutdoors,omitempty" enums:"none,under_30min,30min_1hr,over_1hr"`
	Notes              string `json:"notes,omitempty"`
}

type EnergyLevelsRangeResponse []EnergyLevelsResponse

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
