package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"energyjournal/internal/server"
)

// @title Secret Santa API
// @version 1.0
// @description HTTP API for creating groups, adding gifters, and running a Secret Santa draw.
// @BasePath /
func main() {
	logrus.SetLevel(logrus.DebugLevel)
	// Wire real services into the server here once implementations exist.
	srv := server.New(":8888")
	log.Printf("HTTP server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
