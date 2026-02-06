package server

import (
	"net/http"

	"energyjournal/internal/pkg/httputil"
)

func WriteError(w http.ResponseWriter, err error) {
	httputil.WriteError(w, err)
}
