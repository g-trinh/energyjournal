package calendar

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"energyjournal/internal/domain/calendar"
	errpkg "energyjournal/internal/pkg/error"
)

type fakeRepo struct {
	getFn    func(ctx context.Context, uid string) (*calendar.CalendarConnection, error)
	upsertFn func(ctx context.Context, conn calendar.CalendarConnection) error
}

func (r *fakeRepo) Get(ctx context.Context, uid string) (*calendar.CalendarConnection, error) {
	return r.getFn(ctx, uid)
}

func (r *fakeRepo) Upsert(ctx context.Context, conn calendar.CalendarConnection) error {
	if r.upsertFn != nil {
		return r.upsertFn(ctx, conn)
	}
	return nil
}

type fakeCalendarClient struct {
	calendars []calendar.CalendarItem
	events    []calendar.Event
}

func (c *fakeCalendarClient) ListCalendars(context.Context, string) ([]calendar.CalendarItem, error) {
	return c.calendars, nil
}

func (c *fakeCalendarClient) ListEvents(context.Context, string, string, time.Time, time.Time) ([]calendar.Event, error) {
	return c.events, nil
}

type fakeTokenSource struct {
	token *oauth2.Token
	err   error
}

func (s fakeTokenSource) Token() (*oauth2.Token, error) {
	return s.token, s.err
}

type fakeOAuth struct {
	authURL      string
	exchangeTok  *oauth2.Token
	exchangeErr  error
	tokenSource  oauth2.TokenSource
}

func (o *fakeOAuth) AuthCodeURL(state string, _ ...oauth2.AuthCodeOption) string {
	if o.authURL != "" {
		return o.authURL + "?state=" + state
	}
	return "https://accounts.example/oauth?state=" + state
}

func (o *fakeOAuth) Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return o.exchangeTok, o.exchangeErr
}

func (o *fakeOAuth) TokenSource(context.Context, *oauth2.Token) oauth2.TokenSource {
	return o.tokenSource
}

func TestGetStatus(t *testing.T) {
	t.Parallel()

	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return nil, nil
		},
	}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")

	status, err := svc.GetStatus(context.Background(), "uid")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if status != calendar.StatusDisconnected {
		t.Fatalf("expected disconnected, got %s", status)
	}
}

func TestHandleCallbackInvalidState(t *testing.T) {
	t.Parallel()

	svc := NewCalendarService(&fakeRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")
	err := svc.HandleCallback(context.Background(), "code", "invalid")
	if err == nil {
		t.Fatal("expected error")
	}
	var validationErr *errpkg.InputValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestHandleCallbackUpsertsTokens(t *testing.T) {
	t.Parallel()

	oauth := &fakeOAuth{
		exchangeTok: &oauth2.Token{
			AccessToken:  "access",
			RefreshToken: "refresh",
			Expiry:       time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC),
		},
	}

	var saved calendar.CalendarConnection
	svc := NewCalendarService(&fakeRepo{
		upsertFn: func(_ context.Context, conn calendar.CalendarConnection) error {
			saved = conn
			return nil
		},
	}, &fakeCalendarClient{}, oauth, "secret")
	svc.now = func() time.Time { return time.Date(2026, 3, 3, 11, 0, 0, 0, time.UTC) }

	state := svc.signState("uid-1", svc.now())
	if err := svc.HandleCallback(context.Background(), "code", state); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if saved.UID != "uid-1" || saved.AccessToken != "access" || saved.RefreshToken != "refresh" {
		t.Fatalf("unexpected saved connection: %+v", saved)
	}
}

func TestGetCalendarsRequiresOAuthConnection(t *testing.T) {
	t.Parallel()

	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return nil, nil
		},
	}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")

	_, err := svc.GetCalendars(context.Background(), "uid")
	if err == nil {
		t.Fatal("expected error")
	}
	var connErr *errpkg.CalendarNotConnectedError
	if !errors.As(err, &connErr) {
		t.Fatalf("expected CalendarNotConnectedError, got %T", err)
	}
}

func TestSetCalendarPersistsSelection(t *testing.T) {
	t.Parallel()

	var saved calendar.CalendarConnection
	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return &calendar.CalendarConnection{
				UID:          "uid",
				AccessToken:  "a",
				RefreshToken: "r",
			}, nil
		},
		upsertFn: func(_ context.Context, conn calendar.CalendarConnection) error {
			saved = conn
			return nil
		},
	}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")

	if err := svc.SetCalendar(context.Background(), "uid", "primary"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if saved.CalendarID != "primary" {
		t.Fatalf("expected saved calendarID primary, got %q", saved.CalendarID)
	}
}

func TestGetSpendingAggregatesByColorAndRefreshesToken(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC)
	var upsertCount int
	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return &calendar.CalendarConnection{
				UID:          "uid",
				CalendarID:   "primary",
				AccessToken:  "old",
				RefreshToken: "refresh",
				Expiry:       now.Add(-time.Minute),
			}, nil
		},
		upsertFn: func(_ context.Context, conn calendar.CalendarConnection) error {
			upsertCount++
			if conn.AccessToken != "new-access" {
				t.Fatalf("expected refreshed access token, got %q", conn.AccessToken)
			}
			return nil
		},
	}, &fakeCalendarClient{
		events: []calendar.Event{
			{ColorID: "5", Start: now.Add(-4 * time.Hour), End: now.Add(-3 * time.Hour)},
			{ColorID: "5", Start: now.Add(-3 * time.Hour), End: now.Add(-90 * time.Minute)},
			{ColorID: "1", Start: now.Add(-90 * time.Minute), End: now},
			{ColorID: "", Start: now, End: now}, // skipped
		},
	}, &fakeOAuth{
		tokenSource: fakeTokenSource{
			token: &oauth2.Token{
				AccessToken:  "new-access",
				RefreshToken: "new-refresh",
				Expiry:       now.Add(time.Hour),
			},
		},
	}, "secret")
	svc.now = func() time.Time { return now }

	result, err := svc.GetSpending(context.Background(), "uid", now.Add(-24*time.Hour), now)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if upsertCount != 1 {
		t.Fatalf("expected one upsert after refresh, got %d", upsertCount)
	}
	if result["Sage"] != 2.5 {
		t.Fatalf("expected Sage=2.5, got %v", result["Sage"])
	}
	if result["Tomato"] != 1.5 {
		t.Fatalf("expected Tomato=1.5, got %v", result["Tomato"])
	}
}

func TestVerifyStateExpired(t *testing.T) {
	t.Parallel()

	svc := NewCalendarService(&fakeRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")
	svc.now = func() time.Time { return time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC) }
	expired := svc.signState("uid", svc.now().Add(-20*time.Minute))

	_, err := svc.verifyState(expired)
	if err == nil {
		t.Fatal("expected expiry error")
	}
}
