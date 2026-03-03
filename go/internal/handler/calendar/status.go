package calendar

import (
	"net/http"

	"energyjournal/internal/domain/calendar"
	"energyjournal/internal/server/middleware"
)

type StatusHandler struct {
	service calendar.CalendarService
}

func NewStatusHandler(service calendar.CalendarService) *StatusHandler {
	return &StatusHandler{service: service}
}

// GetStatus godoc
// @Summary Get Google Calendar connection status
// @Description Returns a tri-state status: disconnected (no OAuth), pending_selection (OAuth done, no calendar chosen), connected (ready).
// @Tags calendar
// @Security BearerAuth
// @Success 200 {object} calendar.StatusResponse
// @Failure 401 {object} calendar.ErrorResponse
// @Failure 500 {object} calendar.ErrorResponse
// @Router /calendar/status [get]
func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UIDFromContext(r.Context())
	if !ok || uid == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	status, err := h.service.GetStatus(r.Context(), uid)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, StatusResponse{Status: status})
}
