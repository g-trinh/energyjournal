package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"energyjournal/internal/domain/calendar"
	"energyjournal/internal/domain/user"
	calendarhandler "energyjournal/internal/handler/calendar"
	userhandler "energyjournal/internal/handler/user"
	"energyjournal/internal/pkg/firebase"
	"energyjournal/internal/pkg/firestore"
	"energyjournal/internal/server/middleware"
	calendarservice "energyjournal/internal/service/calendar"
	userservice "energyjournal/internal/service/user"
	userstorage "energyjournal/internal/service/user/storage"
)

// Dependencies groups external services that the HTTP server needs.
type Dependencies struct {
	SpendingService calendar.SpendingService
	UserService     user.UserService
	AuthMiddleware  *middleware.AuthMiddleware
}

// New creates the HTTP server with the default routes.
func New(addr string) *http.Server {
	mux := http.NewServeMux()
	ctx := context.Background()

	firebaseClient, err := firebase.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase client: %v", err)
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable is required")
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}

	userRepo := userstorage.NewUserRepository(firestoreClient.Client)
	tokenRepo := userstorage.NewActivationTokenRepository(firestoreClient.Client)
	authProvider := firebase.NewAuthProvider(firebaseClient, os.Getenv("FIREBASE_API_KEY"))
	emailSender := &noopEmailSender{} // TODO: implement real email sender

	activationBaseURL := lookupEnvOrDefault("FRONTEND_ACTIVATION_BASE_URL", "http://localhost:8080")

	userService := userservice.NewUserService(userRepo, tokenRepo, authProvider, emailSender, activationBaseURL)
	authMiddleware := middleware.NewAuthMiddleware(firebaseClient, userRepo)

	deps := Dependencies{
		SpendingService: calendarservice.NewSpendingService(),
		UserService:     userService,
		AuthMiddleware:  authMiddleware,
	}
	register(mux, deps)

	return &http.Server{
		Addr:    addr,
		Handler: applyCORS(mux),
	}
}

// noopEmailSender is a placeholder email sender that does nothing.
// TODO: Replace with real implementation (e.g., Brevo, SendGrid).
type noopEmailSender struct{}

func (s *noopEmailSender) SendActivationEmail(ctx context.Context, email, activationLink string) error {
	log.Printf("Would send activation email to %s with link %s", email, activationLink)
	return nil
}

// register wires all HTTP handlers onto the given mux.
func register(mux *http.ServeMux, deps Dependencies) {
	mux.HandleFunc("/healthz", health)

	spendingHandler := calendarhandler.NewSpendingHandler(deps.SpendingService)
	if deps.AuthMiddleware != nil {
		mux.Handle("GET /calendar/spending", deps.AuthMiddleware.RequireAuth(http.HandlerFunc(spendingHandler.GetSpending)))
	} else {
		// Fail closed when auth middleware is unavailable.
		mux.HandleFunc("GET /calendar/spending", func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}

	// User routes
	if deps.UserService != nil && deps.AuthMiddleware != nil {
		userHandler := userhandler.NewUserHandler(deps.UserService)

		// POST /users - requires auth but no active status check
		mux.Handle("POST /users", http.HandlerFunc(userHandler.Create))

		// POST /users/activate - no auth required
		NewRoute(mux, http.MethodPost, "/users/activate", userHandler.Activate)

		// POST /users/login - no auth required
		NewRoute(mux, http.MethodPost, "/users/login", userHandler.Login)

		// POST /users/refresh - no auth required
		NewRoute(mux, http.MethodPost, "/users/refresh", userHandler.RefreshToken)

		// Protected routes - require active user
		mux.Handle("GET /users/me", deps.AuthMiddleware.RequireActiveUser(http.HandlerFunc(userHandler.GetProfile)))
		mux.Handle("PUT /users/me", deps.AuthMiddleware.RequireActiveUser(http.HandlerFunc(userHandler.UpdateProfile)))
		mux.Handle("DELETE /users/me", deps.AuthMiddleware.RequireActiveUser(http.HandlerFunc(userHandler.DeleteProfile)))
	}
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
