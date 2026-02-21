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
