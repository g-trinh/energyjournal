package energy

import (
	"encoding/json"
	"errors"
	"net/http"

	"energyjournal/internal/domain/energy"
	pkgerror "energyjournal/internal/pkg/error"
	"energyjournal/internal/server/middleware"
)

type EnergyHandler struct {
	service energy.EnergyService
}

func New(service energy.EnergyService) *EnergyHandler {
	return &EnergyHandler{service: service}
}

// GetLevels godoc
// @Summary Get energy levels for a specific date
// @Tags energy
// @Security BearerAuth
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} energy.EnergyLevelsResponse
// @Failure 400 {object} energy.ErrorResponse
// @Failure 401 {object} energy.ErrorResponse
// @Failure 404 {object} energy.ErrorResponse
// @Failure 500 {object} energy.ErrorResponse
// @Router /energy/levels [get]
func (h *EnergyHandler) GetLevels(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "date query parameter is required"})
		return
	}

	levels, err := h.service.GetByDate(r.Context(), u.UID, date)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, EnergyLevelsResponse{
		Date:      levels.Date,
		Physical:  levels.Physical,
		Mental:    levels.Mental,
		Emotional: levels.Emotional,
	})
}

// GetLevelsByRange godoc
// @Summary Get energy levels for a date range
// @Description Returns all energy levels recorded between from and to (inclusive). Dates with no record are omitted. Max range is 30 days (silently clamped).
// @Tags energy
// @Security BearerAuth
// @Param from query string true "Start date (YYYY-MM-DD, inclusive)" example(2026-02-09)
// @Param to query string true "End date (YYYY-MM-DD, inclusive)" example(2026-02-23)
// @Success 200 {array} energy.EnergyLevelsResponse
// @Failure 400 {object} energy.ErrorResponse
// @Failure 401 {object} energy.ErrorResponse
// @Failure 403 {object} energy.ErrorResponse
// @Failure 500 {object} energy.ErrorResponse
// @Router /energy/levels/range [get]
func (h *EnergyHandler) GetLevelsByRange(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	levels, err := h.service.GetByDateRange(r.Context(), u.UID, from, to)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	response := make([]EnergyLevelsResponse, 0, len(levels))
	for _, level := range levels {
		response = append(response, EnergyLevelsResponse{
			Date:      level.Date,
			Physical:  level.Physical,
			Mental:    level.Mental,
			Emotional: level.Emotional,
		})
	}

	writeJSON(w, http.StatusOK, EnergyLevelsRangeResponse(response))
}

// SaveLevels godoc
// @Summary Save energy levels for a specific date
// @Tags energy
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body energy.SaveEnergyLevelsRequest true "Energy levels data"
// @Success 200 {object} energy.EnergyLevelsResponse
// @Failure 400 {object} energy.ErrorResponse
// @Failure 401 {object} energy.ErrorResponse
// @Failure 403 {object} energy.ErrorResponse
// @Failure 500 {object} energy.ErrorResponse
// @Router /energy/levels [put]
func (h *EnergyHandler) SaveLevels(w http.ResponseWriter, r *http.Request) {
	u, ok := middleware.UserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var req SaveEnergyLevelsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	levels := energy.EnergyLevels{
		UID:       u.UID,
		Date:      req.Date,
		Physical:  req.Physical,
		Mental:    req.Mental,
		Emotional: req.Emotional,
	}
	if err := h.service.Save(r.Context(), levels); err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, EnergyLevelsResponse{
		Date:      req.Date,
		Physical:  req.Physical,
		Mental:    req.Mental,
		Emotional: req.Emotional,
	})
}

func writeDomainError(w http.ResponseWriter, err error) {
	var validationErr *pkgerror.InputValidationError
	if errors.As(err, &validationErr) {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: validationErr.Error()})
		return
	}

	var notFoundErr *pkgerror.NotFoundError
	if errors.As(err, &notFoundErr) {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: notFoundErr.Error()})
		return
	}

	writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
}
