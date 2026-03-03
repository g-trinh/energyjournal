package calendar

import (
	"encoding/json"
	"errors"
	"net/http"

	errpkg "energyjournal/internal/pkg/error"
	"energyjournal/internal/pkg/httputil"
)

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	var notConnected *errpkg.CalendarNotConnectedError
	if errors.As(err, &notConnected) {
		writeJSON(w, http.StatusFailedDependency, ErrorResponse{Error: notConnected.Error()})
		return
	}

	statusCode, message := httputil.MapErrors(err)
	writeJSON(w, statusCode, ErrorResponse{Error: message})
}
