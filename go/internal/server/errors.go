package server

import (
	"encoding/json"
	"errors"
	"net/http"

	errpkg "energyjournal/internal/pkg/error"
	"energyjournal/internal/pkg/httputil"
)

func WriteError(w http.ResponseWriter, err error) {
	var calendarErr *errpkg.CalendarNotConnectedError
	if errors.As(err, &calendarErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusFailedDependency)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": calendarErr.Error()})
		return
	}

	httputil.WriteError(w, err)
}
