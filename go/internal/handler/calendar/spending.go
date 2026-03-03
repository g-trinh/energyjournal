package calendar

import (
	"encoding/json"
	"net/http"
	"time"

	"energyjournal/internal/domain/calendar"
	errpkg "energyjournal/internal/pkg/error"
	"energyjournal/internal/server/middleware"
)

const dateFormat = "2006-01-02"

// SpendingHandler handles HTTP requests for calendar spending.
type SpendingHandler struct {
	service calendar.CalendarService
}

// NewSpendingHandler creates a new SpendingHandler.
func NewSpendingHandler(service calendar.CalendarService) *SpendingHandler {
	return &SpendingHandler{service: service}
}

// GetSpending handles GET /calendar/spending requests.
// @Summary Get time spendings from the selected Google Calendar
// @Description Aggregates event durations from the user's selected Google Calendar grouped by event color label.
// @Tags calendar
// @Security BearerAuth
// @Param start query string true "Start date (YYYY-MM-DD)"
// @Param end query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} calendar.Spendings
// @Failure 400 {object} calendar.ErrorResponse
// @Failure 401 {object} calendar.ErrorResponse
// @Failure 424 {object} calendar.ErrorResponse
// @Failure 500 {object} calendar.ErrorResponse
// @Router /calendar/spending [get]
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

	uid, ok := middleware.UIDFromContext(r.Context())
	if !ok || uid == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	spendings, err := h.service.GetSpending(r.Context(), uid, start, end)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(spendings)
}
