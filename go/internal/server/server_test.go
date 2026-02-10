package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"energyjournal/internal/domain/calendar"
	"energyjournal/internal/server/middleware"
	"firebase.google.com/go/v4/auth"
)

type stubSpendingService struct {
	getSpending func(start, end time.Time) (calendar.Spendings, error)
}

func (s *stubSpendingService) GetSpending(start, end time.Time) (calendar.Spendings, error) {
	return s.getSpending(start, end)
}

type stubVerifier struct {
	verifyIDToken func(ctx context.Context, idToken string) (*auth.Token, error)
}

func (v *stubVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return v.verifyIDToken(ctx, idToken)
}

func TestRegister_CalendarSpending_UnauthorizedWithoutToken(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	register(mux, Dependencies{
		SpendingService: &stubSpendingService{
			getSpending: func(start, end time.Time) (calendar.Spendings, error) {
				t.Fatal("spending service should not be called")
				return nil, nil
			},
		},
		AuthMiddleware: middleware.NewAuthMiddlewareWithVerifier(&stubVerifier{
			verifyIDToken: func(ctx context.Context, idToken string) (*auth.Token, error) {
				t.Fatal("token verifier should not be called")
				return nil, nil
			},
		}, nil),
	})

	req := httptest.NewRequest(http.MethodGet, "/calendar/spending?start=2026-01-01&end=2026-01-31", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRegister_CalendarSpending_AuthorizedWithToken(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	register(mux, Dependencies{
		SpendingService: &stubSpendingService{
			getSpending: func(start, end time.Time) (calendar.Spendings, error) {
				return calendar.Spendings{"Work": 12.5}, nil
			},
		},
		AuthMiddleware: middleware.NewAuthMiddlewareWithVerifier(&stubVerifier{
			verifyIDToken: func(ctx context.Context, idToken string) (*auth.Token, error) {
				return &auth.Token{
					UID:    "uid-1",
					Claims: map[string]any{"email": "user@example.com"},
				}, nil
			},
		}, nil),
	})

	req := httptest.NewRequest(http.MethodGet, "/calendar/spending?start=2026-01-01&end=2026-01-31", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got map[string]float64
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if got["Work"] != 12.5 {
		t.Fatalf("expected Work to be 12.5, got %v", got["Work"])
	}
}

func TestRegister_CalendarSpending_WithoutAuthMiddlewareFailsClosed(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	register(mux, Dependencies{
		SpendingService: &stubSpendingService{
			getSpending: func(start, end time.Time) (calendar.Spendings, error) {
				t.Fatal("spending service should not be called")
				return nil, nil
			},
		},
		AuthMiddleware: nil,
	})

	req := httptest.NewRequest(http.MethodGet, "/calendar/spending?start=2026-01-01&end=2026-01-31", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
