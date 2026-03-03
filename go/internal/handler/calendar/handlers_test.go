package calendar

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"energyjournal/internal/domain/calendar"
	"energyjournal/internal/server/middleware"
)

type stubCalendarService struct {
	getStatus    func(ctx context.Context, uid string) (calendar.ConnectionStatus, error)
	buildAuthURL func(uid string) string
	callback     func(ctx context.Context, code, state string) error
	getCalendars func(ctx context.Context, uid string) ([]calendar.CalendarItem, error)
	setCalendar  func(ctx context.Context, uid, calendarID string) error
	getSpending  func(ctx context.Context, uid string, start, end time.Time) (calendar.Spendings, error)
}

func (s *stubCalendarService) GetStatus(ctx context.Context, uid string) (calendar.ConnectionStatus, error) {
	return s.getStatus(ctx, uid)
}
func (s *stubCalendarService) BuildAuthURL(uid string) string {
	return s.buildAuthURL(uid)
}
func (s *stubCalendarService) HandleCallback(ctx context.Context, code, state string) error {
	return s.callback(ctx, code, state)
}
func (s *stubCalendarService) GetCalendars(ctx context.Context, uid string) ([]calendar.CalendarItem, error) {
	return s.getCalendars(ctx, uid)
}
func (s *stubCalendarService) SetCalendar(ctx context.Context, uid, calendarID string) error {
	return s.setCalendar(ctx, uid, calendarID)
}
func (s *stubCalendarService) GetSpending(ctx context.Context, uid string, start, end time.Time) (calendar.Spendings, error) {
	return s.getSpending(ctx, uid, start, end)
}

func TestStatusHandlerResponseShape(t *testing.T) {
	t.Parallel()

	handler := NewCalendarHandler(&stubCalendarService{
		getStatus: func(ctx context.Context, uid string) (calendar.ConnectionStatus, error) {
			return calendar.StatusConnected, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/calendar/status", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeyUID, "uid-1"))
	rr := httptest.NewRecorder()
	handler.GetStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var payload StatusResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Status != calendar.StatusConnected {
		t.Fatalf("unexpected status: %s", payload.Status)
	}
}

func TestOAuthGetAuthURLResponseShape(t *testing.T) {
	t.Parallel()

	handler := NewOAuthHandler(&stubCalendarService{
		buildAuthURL: func(uid string) string {
			return "https://accounts.google.com/o/oauth2/auth?state=abc"
		},
	}, "http://localhost:8080")

	req := httptest.NewRequest(http.MethodGet, "/calendar/auth", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeyUID, "uid-1"))
	rr := httptest.NewRecorder()
	handler.GetAuthURL(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var payload AuthURLResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.AuthURL == "" {
		t.Fatal("expected non-empty auth_url")
	}
}

func TestCalendarsSetConnectionParsesRequest(t *testing.T) {
	t.Parallel()

	captured := ""
	handler := NewCalendarHandler(&stubCalendarService{
		setCalendar: func(ctx context.Context, uid, calendarID string) error {
			captured = calendarID
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodPut, "/calendar/connection", strings.NewReader(`{"calendar_id":"primary"}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeyUID, "uid-1"))
	rr := httptest.NewRecorder()
	handler.SetConnection(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if captured != "primary" {
		t.Fatalf("expected captured calendar id primary, got %q", captured)
	}
}

func TestSpendingHandlerParsesDatesAndReturnsMap(t *testing.T) {
	t.Parallel()

	handler := NewSpendingHandler(&stubCalendarService{
		getSpending: func(ctx context.Context, uid string, start, end time.Time) (calendar.Spendings, error) {
			if uid != "uid-1" {
				t.Fatalf("unexpected uid %s", uid)
			}
			return calendar.Spendings{"Sage": 3.5}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/calendar/spending?start=2026-03-01&end=2026-03-03", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeyUID, "uid-1"))
	rr := httptest.NewRecorder()
	handler.GetSpending(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var payload map[string]float64
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["Sage"] != 3.5 {
		t.Fatalf("expected Sage=3.5, got %v", payload["Sage"])
	}
}
