package calendar

import (
	"encoding/json"
	"net/http"

	"energyjournal/internal/domain/calendar"
	errpkg "energyjournal/internal/pkg/error"
	"energyjournal/internal/server/middleware"
)

type CalendarsHandler struct {
	service calendar.CalendarService
}

func NewCalendarsHandler(service calendar.CalendarService) *CalendarsHandler {
	return &CalendarsHandler{service: service}
}

func (h *CalendarsHandler) GetCalendars(w http.ResponseWriter, r *http.Request) {
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

func (h *CalendarsHandler) SetConnection(w http.ResponseWriter, r *http.Request) {
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
