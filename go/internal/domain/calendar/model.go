package calendar

import (
	"context"
	"time"
)

// Spendings represents time spent per event type (in hours).
type Spendings map[string]float64

type ConnectionStatus string

const (
	StatusDisconnected     ConnectionStatus = "disconnected"
	StatusPendingSelection ConnectionStatus = "pending_selection"
	StatusConnected        ConnectionStatus = "connected"
)

type CalendarConnection struct {
	UID          string
	CalendarID   string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

type CalendarItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Event struct {
	ColorID string
	Start   time.Time
	End     time.Time
}

type CalendarConnectionRepository interface {
	Get(ctx context.Context, uid string) (*CalendarConnection, error)
	Upsert(ctx context.Context, conn CalendarConnection) error
}

// SpendingService defines the contract for spending retrieval.
// Kept for backward compatibility while handlers migrate to CalendarService.
type SpendingService interface {
	GetSpending(start, end time.Time) (Spendings, error)
}

// CalendarService defines the full Google Calendar feature contract.
type CalendarService interface {
	GetStatus(ctx context.Context, uid string) (ConnectionStatus, error)
	BuildAuthURL(uid string) string
	HandleCallback(ctx context.Context, code, state string) error
	GetCalendars(ctx context.Context, uid string) ([]CalendarItem, error)
	SetCalendar(ctx context.Context, uid, calendarID string) error
	GetSpending(ctx context.Context, uid string, start, end time.Time) (Spendings, error)
}
