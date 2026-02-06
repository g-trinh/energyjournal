package httputil

import (
	"errors"
	"net/http"

	errpkg "energyjournal/internal/pkg/error"
)

func MapErrors(err error) (statusCode int, message string) {
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

func WriteError(w http.ResponseWriter, err error) {
	statusCode, message := MapErrors(err)
	http.Error(w, message, statusCode)
}
