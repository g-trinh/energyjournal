package calendar

import (
	"encoding/json"
	"net/http"

	"energyjournal/internal/domain/calendar"
	errpkg "energyjournal/internal/pkg/error"
	"energyjournal/internal/server/middleware"
)

type CalendarHandler struct {
	service calendar.CalendarService
}

func NewCalendarHandler(service calendar.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
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
func (h *CalendarHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
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

// GetCalendars godoc
// @Summary List user's Google Calendars
// @Tags calendar
// @Security BearerAuth
// @Success 200 {array} calendar.CalendarItem
// @Failure 401 {object} calendar.ErrorResponse
// @Failure 424 {object} calendar.ErrorResponse
// @Failure 500 {object} calendar.ErrorResponse
// @Router /calendar/calendars [get]
func (h *CalendarHandler) GetCalendars(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UIDFromContext(r.Context())
	if !ok || uid == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	calendars, err := h.service.GetCalendars(r.Context(), uid)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, calendars)
}

// SetConnection godoc
// @Summary Save selected calendar
// @Description Persists the user's chosen Google Calendar ID.
// @Tags calendar
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body calendar.SetConnectionRequest true "Selected calendar ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} calendar.ErrorResponse
// @Failure 401 {object} calendar.ErrorResponse
// @Failure 424 {object} calendar.ErrorResponse
// @Failure 500 {object} calendar.ErrorResponse
// @Router /calendar/connection [put]
func (h *CalendarHandler) SetConnection(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UIDFromContext(r.Context())
	if !ok || uid == "" {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req SetConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}
	if req.CalendarID == "" {
		writeError(w, errpkg.NewInputValidationError("calendar_id", "required"))
		return
	}

	if err := h.service.SetCalendar(r.Context(), uid, req.CalendarID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
