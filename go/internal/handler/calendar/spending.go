package calendar

import (
	"encoding/json"
	"net/http"
	"time"

	"energyjournal/internal/domain/calendar"
	errpkg "energyjournal/internal/pkg/error"
)

const dateFormat = "2006-01-02"

// SpendingHandler handles HTTP requests for calendar spending.
type SpendingHandler struct {
	service calendar.SpendingService
}

// NewSpendingHandler creates a new SpendingHandler.
func NewSpendingHandler(service calendar.SpendingService) *SpendingHandler {
	return &SpendingHandler{service: service}
}

// GetSpending handles GET /calendar/spending requests.
func (h *SpendingHandler) GetSpending(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	if startStr == "" {
		writeError(w, errpkg.NewInputValidationError("start", "required"))
		return
	}

	endStr := r.URL.Query().Get("end")
	if endStr == "" {
		writeError(w, errpkg.NewInputValidationError("end", "required"))
		return
	}

	start, err := time.Parse(dateFormat, startStr)
	if err != nil {
		writeError(w, errpkg.NewInputValidationError("start", "invalid date format, expected YYYY-MM-DD"))
		return
	}

	end, err := time.Parse(dateFormat, endStr)
	if err != nil {
		writeError(w, errpkg.NewInputValidationError("end", "invalid date format, expected YYYY-MM-DD"))
		return
	}

	spendings, err := h.service.GetSpending(start, end)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spendings)
}

func writeError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	switch e := err.(type) {
	case *errpkg.InputValidationError:
		statusCode = http.StatusBadRequest
		message = e.Error()
	case *errpkg.NotFoundError:
		statusCode = http.StatusNotFound
		message = e.Error()
	default:
		statusCode = http.StatusInternalServerError
		message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
