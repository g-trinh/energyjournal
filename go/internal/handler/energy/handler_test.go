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
				UID:                uid,
				Date:               date,
				Physical:           7,
				Mental:             5,
				Emotional:          8,
				SleepQuality:       4,
				StressLevel:        2,
				PhysicalActivity:   "light",
				Nutrition:          "good",
				SocialInteractions: "positive",
				TimeOutdoors:       "30min_1hr",
				Notes:              "Felt balanced.",
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
	if payload.SleepQuality != 4 || payload.StressLevel != 2 || payload.PhysicalActivity != "light" {
		t.Fatalf("unexpected context payload: %+v", payload)
	}
	if payload.Nutrition != "good" || payload.SocialInteractions != "positive" || payload.TimeOutdoors != "30min_1hr" {
		t.Fatalf("unexpected context payload: %+v", payload)
	}
	if payload.Notes != "Felt balanced." {
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

func TestEnergyHandler_GetLevelsByRange_MissingFromDelegatesToService(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if from != "" || to != "2026-02-21" {
				t.Fatalf("unexpected range values: from=%q to=%q", from, to)
			}
			return []energy.EnergyLevels{}, nil
		},
	})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels/range?to=2026-02-21", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevelsByRange(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestEnergyHandler_GetLevelsByRange_MissingToDelegatesToService(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if from != "2026-02-01" || to != "" {
				t.Fatalf("unexpected range values: from=%q to=%q", from, to)
			}
			return []energy.EnergyLevels{}, nil
		},
	})
	req := withUserContext(httptest.NewRequest(http.MethodGet, "/energy/levels/range?from=2026-02-01", nil), "uid-1")
	rr := httptest.NewRecorder()

	handler.GetLevelsByRange(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestEnergyHandler_GetLevelsByRange_Success(t *testing.T) {
	t.Parallel()

	handler := New(&stubEnergyService{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			return []energy.EnergyLevels{
				{
					UID:                uid,
					Date:               "2026-02-20",
					Physical:           6,
					Mental:             5,
					Emotional:          4,
					SleepQuality:       3,
					StressLevel:        1,
					PhysicalActivity:   "none",
					Nutrition:          "average",
					SocialInteractions: "neutral",
					TimeOutdoors:       "under_30min",
					Notes:              "Rushed day.",
				},
				{
					UID:                uid,
					Date:               "2026-02-21",
					Physical:           7,
					Mental:             6,
					Emotional:          8,
					SleepQuality:       5,
					StressLevel:        2,
					PhysicalActivity:   "moderate",
					Nutrition:          "good",
					SocialInteractions: "positive",
					TimeOutdoors:       "30min_1hr",
					Notes:              "Great momentum.",
				},
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
	if payload[0].PhysicalActivity != "none" || payload[1].PhysicalActivity != "moderate" {
		t.Fatalf("expected context fields in range response, got %+v", payload)
	}
	if payload[0].SleepQuality != 3 || payload[1].SleepQuality != 5 {
		t.Fatalf("expected context fields in range response, got %+v", payload)
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
			if levels.SleepQuality != 4 || levels.StressLevel != 3 {
				t.Fatalf("expected context int fields to be mapped, got %+v", levels)
			}
			if levels.PhysicalActivity != "light" || levels.Nutrition != "good" {
				t.Fatalf("expected context enum fields to be mapped, got %+v", levels)
			}
			if levels.SocialInteractions != "positive" || levels.TimeOutdoors != "over_1hr" {
				t.Fatalf("expected context enum fields to be mapped, got %+v", levels)
			}
			if levels.Notes != "Clear focus today." {
				t.Fatalf("expected notes to be mapped, got %+v", levels)
			}
			return nil
		},
	})
	body := bytes.NewBufferString(`{"date":"2026-02-21","physical":7,"mental":5,"emotional":8,"sleepQuality":4,"stressLevel":3,"physicalActivity":"light","nutrition":"good","socialInteractions":"positive","timeOutdoors":"over_1hr","notes":"Clear focus today."}`)
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
	if payload.SleepQuality != 4 || payload.StressLevel != 3 || payload.PhysicalActivity != "light" {
		t.Fatalf("unexpected context payload: %+v", payload)
	}
	if payload.Nutrition != "good" || payload.SocialInteractions != "positive" || payload.TimeOutdoors != "over_1hr" {
		t.Fatalf("unexpected context payload: %+v", payload)
	}
	if payload.Notes != "Clear focus today." {
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
