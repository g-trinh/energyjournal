package energy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"energyjournal/internal/domain/energy"
	"energyjournal/internal/domain/user"
	pkgerror "energyjournal/internal/pkg/error"
	"energyjournal/internal/server/middleware"
)

func TestEnergyHandler_GetLevels_MissingDateReturnsBadRequest(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevels(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEnergyHandler_GetLevels_Success(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		getByDate: func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
			return &energy.EnergyLevels{
				UID:       uid,
				Date:      date,
				Physical:  7,
				Mental:    5,
				Emotional: 8,
			}, nil
		},
	})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels?date=2026-02-21", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevels(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var payload EnergyLevelsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Date != "2026-02-21" || payload.Physical != 7 || payload.Mental != 5 || payload.Emotional != 8 {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestEnergyHandler_GetLevels_NotFoundReturns404(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		getByDate: func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
			return nil, pkgerror.NewNotFoundError("energy_levels", "uid-1_2026-02-21")
		},
	})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels?date=2026-02-21", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevels(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}

func TestEnergyHandler_GetLevelsByRange_MissingFromReturnsBadRequest(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels/range?to=2026-02-21", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevelsByRange(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEnergyHandler_GetLevelsByRange_MissingToReturnsBadRequest(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels/range?from=2026-02-01", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevelsByRange(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEnergyHandler_GetLevelsByRange_Success(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			return []energy.EnergyLevels{
				{UID: uid, Date: "2026-02-20", Physical: 6, Mental: 5, Emotional: 4},
				{UID: uid, Date: "2026-02-21", Physical: 7, Mental: 6, Emotional: 8},
			}, nil
		},
	})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels/range?from=2026-02-20&to=2026-02-21", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevelsByRange(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var payload []EnergyLevelsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(payload) != 2 || payload[0].Date != "2026-02-20" || payload[1].Date != "2026-02-21" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestEnergyHandler_SaveLevels_MalformedBodyReturnsBadRequest(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{})
	req := withUserContext(httptest.NewRequest(http.MethodPut, "/energy/levels", bytes.NewBufferString("{")), "uid-1")
	rr := httptest.NewRecorder()

	handler.SaveLevels(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEnergyHandler_SaveLevels_OutOfRangeReturnsBadRequest(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		save: func(ctx context.Context, levels energy.EnergyLevels) error {
			return pkgerror.NewInputValidationError("physical", "must be between 0 and 10")
		},
	})
	body := bytes.NewBufferString(`{"date":"2026-02-21","physical":11,"mental":5,"emotional":8}`)
	req := withUserContext(httptest.NewRequest(http.MethodPut, "/energy/levels", body), "uid-1")
	rr := httptest.NewRecorder()

	handler.SaveLevels(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEnergyHandler_SaveLevels_Success(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		save: func(ctx context.Context, levels energy.EnergyLevels) error {
			return nil
		},
	})
	body := bytes.NewBufferString(`{"date":"2026-02-21","physical":7,"mental":5,"emotional":8}`)
	req := withUserContext(httptest.NewRequest(http.MethodPut, "/energy/levels", body), "uid-1")
	rr := httptest.NewRecorder()

	handler.SaveLevels(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var payload EnergyLevelsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Date != "2026-02-21" || payload.Physical != 7 || payload.Mental != 5 || payload.Emotional != 8 {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestEnergyHandler_SaveLevels_InternalErrorReturns500(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		save: func(ctx context.Context, levels energy.EnergyLevels) error {
			return errors.New("unexpected failure")
		},
	})
	body := bytes.NewBufferString(`{"date":"2026-02-21","physical":7,"mental":5,"emotional":8}`)
	req := withUserContext(httptest.NewRequest(http.MethodPut, "/energy/levels", body), "uid-1")
	rr := httptest.NewRecorder()

	handler.SaveLevels(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

type stubEnergyService struct {
	getByDate      func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error)
	getByDateRange func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error)
	save           func(ctx context.Context, levels energy.EnergyLevels) error
}

func (s *stubEnergyService) GetByDate(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
	if s.getByDate != nil {
		return s.getByDate(ctx, uid, date)
	}
	return nil, nil
}

func (s *stubEnergyService) GetByDateRange(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
	if s.getByDateRange != nil {
		return s.getByDateRange(ctx, uid, from, to)
	}
	return nil, nil
}

func (s *stubEnergyService) Save(ctx context.Context, levels energy.EnergyLevels) error {
	if s.save != nil {
		return s.save(ctx, levels)
	}
	return nil
}

func withUserContext(req *http.Request, uid string) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.ContextKeyUser, &user.User{
		UID:    uid,
		Email:  "user@example.com",
		Status: user.StatusActive,
	})
	return req.WithContext(ctx)
}
