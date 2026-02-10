package middleware

import (
	"context"
	"net/http"
	"strings"

	"energyjournal/internal/domain/user"
	"energyjournal/internal/pkg/firebase"
	"firebase.google.com/go/v4/auth"
)

type ContextKey string

const ContextKeyUser ContextKey = "user"
const ContextKeyUID ContextKey = "uid"
const ContextKeyEmail ContextKey = "email"

type AuthMiddleware struct {
	tokenVerifier IDTokenVerifier
	userRepo      user.UserRepository
}

// IDTokenVerifier validates and decodes ID tokens.
type IDTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

func NewAuthMiddleware(firebaseClient *firebase.Client, userRepo user.UserRepository) *AuthMiddleware {
	return NewAuthMiddlewareWithVerifier(firebaseClient, userRepo)
}

// NewAuthMiddlewareWithVerifier allows injecting a verifier (useful for tests).
func NewAuthMiddlewareWithVerifier(verifier IDTokenVerifier, userRepo user.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenVerifier: verifier,
		userRepo:      userRepo,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if m.tokenVerifier == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		decoded, err := m.tokenVerifier.VerifyIDToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUID, decoded.UID)
		if email, ok := decoded.Claims["email"].(string); ok {
			ctx = context.WithValue(ctx, ContextKeyEmail, email)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireActiveUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if m.tokenVerifier == nil || m.userRepo == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		decoded, err := m.tokenVerifier.VerifyIDToken(r.Context(), token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u, err := m.userRepo.GetByUID(r.Context(), decoded.UID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if u.Status != user.StatusActive {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) (*user.User, bool) {
	u, ok := ctx.Value(ContextKeyUser).(*user.User)
	return u, ok
}

func UIDFromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(ContextKeyUID).(string)
	return uid, ok
}

func EmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyEmail).(string)
	return email, ok
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http.ErrNoCookie
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", http.ErrNoCookie
	}

	return parts[1], nil
}
