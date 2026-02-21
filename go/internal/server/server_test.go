package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"energyjournal/internal/domain/calendar"
	"energyjournal/internal/domain/energy"
	"energyjournal/internal/domain/user"
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

type stubEnergyService struct {
	getByDate func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error)
	save      func(ctx context.Context, levels energy.EnergyLevels) error
}

func (s *stubEnergyService) GetByDate(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
	if s.getByDate != nil {
		return s.getByDate(ctx, uid, date)
	}
	return nil, nil
}

func (s *stubEnergyService) Save(ctx context.Context, levels energy.EnergyLevels) error {
	if s.save != nil {
		return s.save(ctx, levels)
	}
	return nil
}

type stubUserRepo struct {
	getByUID func(ctx context.Context, uid string) (*user.User, error)
}

func (s *stubUserRepo) Create(ctx context.Context, user *user.User) error {
	return nil
}

func (s *stubUserRepo) GetByUID(ctx context.Context, uid string) (*user.User, error) {
	if s.getByUID != nil {
		return s.getByUID(ctx, uid)
	}
	return &user.User{UID: uid, Status: user.StatusActive}, nil
}

func (s *stubUserRepo) Update(ctx context.Context, user *user.User) error {
	return nil
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

func TestRegister_EnergyLevels_UnauthorizedWithoutToken_GET(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	register(mux, Dependencies{
		EnergyService: &stubEnergyService{
			getByDate: func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
				t.Fatal("energy service should not be called")
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

	req := httptest.NewRequest(http.MethodGet, "/energy/levels?date=2026-02-21", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestRegister_EnergyLevels_UnauthorizedWithoutToken_PUT(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	register(mux, Dependencies{
		EnergyService: &stubEnergyService{
			save: func(ctx context.Context, levels energy.EnergyLevels) error {
				t.Fatal("energy service should not be called")
				return nil
			},
		},
		AuthMiddleware: middleware.NewAuthMiddlewareWithVerifier(&stubVerifier{
			verifyIDToken: func(ctx context.Context, idToken string) (*auth.Token, error) {
				t.Fatal("token verifier should not be called")
				return nil, nil
			},
		}, nil),
	})

	req := httptest.NewRequest(http.MethodPut, "/energy/levels", strings.NewReader(`{"date":"2026-02-21","physical":7,"mental":5,"emotional":8}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}
