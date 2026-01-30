package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"energyjournal/internal/domain/calendar"
	calendarhandler "energyjournal/internal/handler/calendar"
	calendarservice "energyjournal/internal/service/calendar"
)

// Dependencies groups external services that the HTTP server needs.
type Dependencies struct {
	SpendingService calendar.SpendingService
}

// New creates the HTTP server with the default routes.
func New(addr string) *http.Server {
	mux := http.NewServeMux()

	deps := Dependencies{
		SpendingService: calendarservice.NewSpendingService(),
	}
	register(mux, deps)

	return &http.Server{
		Addr:    addr,
		Handler: applyCORS(mux),
	}
}

// register wires all HTTP handlers onto the given mux.
func register(mux *http.ServeMux, deps Dependencies) {
	mux.HandleFunc("/healthz", health)

	spendingHandler := calendarhandler.NewSpendingHandler(deps.SpendingService)
	NewRoute(mux, http.MethodGet, "/calendar/spending", spendingHandler.GetSpending)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func NewRoute(mux *http.ServeMux, method string, path string, next http.HandlerFunc) {
	mux.HandleFunc(fmt.Sprintf("%s %s", method, path), next)
}

func applyCORS(next http.Handler) http.Handler {
	allowedOrigins := loadAllowedOrigins()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Add("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loadAllowedOrigins() map[string]bool {
	value := os.Getenv("ALLOWED_ORIGINS")
	if value == "" {
		value = "http://localhost:8080"
	}

	origins := map[string]bool{}
	for _, origin := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		origins[trimmed] = true
	}

	return origins
}

func lookupEnvOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
