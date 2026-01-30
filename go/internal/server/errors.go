package server

import (
	"errors"
	"net/http"

	errpkg "energyjournal/internal/pkg/error"
)

func mapErrors(err error) (statusCode int, message string) {
	var validationErr *errpkg.InputValidationError
	if errors.As(err, &validationErr) {
		return http.StatusBadRequest, validationErr.Error()
	}

	var notFoundErr *errpkg.NotFoundError
	if errors.As(err, &notFoundErr) {
		return http.StatusNotFound, notFoundErr.Error()
	}

	return http.StatusInternalServerError, err.Error()
}
