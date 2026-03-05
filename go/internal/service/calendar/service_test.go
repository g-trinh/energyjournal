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
	calendars      []calendar.CalendarItem
	events         []calendar.Event
	listEventsFn   func(ctx context.Context, token, calendarID string, start, end time.Time) ([]calendar.Event, error)
	listEventsCall int
}

func (c *fakeCalendarClient) ListCalendars(context.Context, string) ([]calendar.CalendarItem, error) {
	return c.calendars, nil
}

func (c *fakeCalendarClient) ListEvents(ctx context.Context, token, calendarID string, start, end time.Time) ([]calendar.Event, error) {
	c.listEventsCall++
	if c.listEventsFn != nil {
		return c.listEventsFn(ctx, token, calendarID, start, end)
	}
	return c.events, nil
}

type fakeSpendingCacheRepo struct {
	getFn func(ctx context.Context, uid string, weekStart time.Time) (calendar.Spendings, error)
	setFn func(ctx context.Context, uid string, weekStart time.Time, spendings calendar.Spendings) error
}

func (r *fakeSpendingCacheRepo) Get(ctx context.Context, uid string, weekStart time.Time) (calendar.Spendings, error) {
	if r.getFn == nil {
		return nil, nil
	}
	return r.getFn(ctx, uid, weekStart)
}

func (r *fakeSpendingCacheRepo) Set(ctx context.Context, uid string, weekStart time.Time, spendings calendar.Spendings) error {
	if r.setFn == nil {
		return nil
	}
	return r.setFn(ctx, uid, weekStart, spendings)
}

type fakeTokenSource struct {
	token *oauth2.Token
	err   error
}

func (s fakeTokenSource) Token() (*oauth2.Token, error) {
	return s.token, s.err
}

type fakeOAuth struct {
	authURL     string
	exchangeTok *oauth2.Token
	exchangeErr error
	tokenSource oauth2.TokenSource
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
	}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")

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

	svc := NewCalendarService(&fakeRepo{}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")
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
	}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{}, oauth, "secret")
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
	}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")

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
	}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")

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
	}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{
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

	svc := NewCalendarService(&fakeRepo{}, &fakeSpendingCacheRepo{}, &fakeCalendarClient{}, &fakeOAuth{}, "secret")
	svc.now = func() time.Time { return time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC) }
	expired := svc.signState("uid", svc.now().Add(-20*time.Minute))

	_, err := svc.verifyState(expired)
	if err == nil {
		t.Fatal("expected expiry error")
	}
}

func TestGetSpendingCacheHitSkipsGoogle(BT *testing.T) {
	BT.Parallel()

	now := time.Date(2026, 3, 5, 12, 0, 0, 0, time.UTC)
	client := &fakeCalendarClient{}
	cacheRepo := &fakeSpendingCacheRepo{
		getFn: func(context.Context, string, time.Time) (calendar.Spendings, error) {
			return calendar.Spendings{"Tomato": 3.5}, nil
		},
	}

	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return &calendar.CalendarConnection{
				UID:          "uid",
				CalendarID:   "primary",
				AccessToken:  "token",
				RefreshToken: "refresh",
				Expiry:       now.Add(time.Hour),
			}, nil
		},
	}, cacheRepo, client, &fakeOAuth{}, "secret")

	result, err := svc.GetSpending(context.Background(), "uid", now, now.AddDate(0, 0, 6))
	if err != nil {
		BT.Fatalf("unexpected err: %v", err)
	}
	if client.listEventsCall != 0 {
		BT.Fatalf("expected no google calls on cache hit, got %d", client.listEventsCall)
	}
	if result["Tomato"] != 3.5 {
		BT.Fatalf("expected cached Tomato=3.5, got %v", result["Tomato"])
	}
}

func TestGetSpendingCacheMissFetchesAndStores(BT *testing.T) {
	BT.Parallel()

	now := time.Date(2026, 3, 5, 12, 0, 0, 0, time.UTC)
	start := time.Date(2026, 3, 3, 0, 0, 0, 0, time.UTC) // Tuesday, validates Monday normalization.
	var stored calendar.Spendings
	var storedWeekStart time.Time

	cacheRepo := &fakeSpendingCacheRepo{
		getFn: func(context.Context, string, time.Time) (calendar.Spendings, error) {
			return nil, nil
		},
		setFn: func(_ context.Context, _ string, weekStart time.Time, spendings calendar.Spendings) error {
			storedWeekStart = weekStart
			stored = spendings
			return nil
		},
	}

	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return &calendar.CalendarConnection{
				UID:          "uid",
				CalendarID:   "primary",
				AccessToken:  "token",
				RefreshToken: "refresh",
				Expiry:       now.Add(time.Hour),
			}, nil
		},
	}, cacheRepo, &fakeCalendarClient{
		events: []calendar.Event{
			{ColorID: "5", Start: now.Add(-3 * time.Hour), End: now.Add(-time.Hour)},
			{ColorID: "999", Start: now.Add(-time.Hour), End: now},
		},
	}, &fakeOAuth{}, "secret")

	result, err := svc.GetSpending(context.Background(), "uid", start, start.AddDate(0, 0, 6))
	if err != nil {
		BT.Fatalf("unexpected err: %v", err)
	}
	expectedMonday := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	if storedWeekStart.Format("2006-01-02") != expectedMonday.Format("2006-01-02") {
		BT.Fatalf("expected weekStart %s, got %s", expectedMonday.Format("2006-01-02"), storedWeekStart.Format("2006-01-02"))
	}
	if result["Sage"] != 2 {
		BT.Fatalf("expected Sage=2, got %v", result["Sage"])
	}
	if result["Default"] != 1 {
		BT.Fatalf("expected Default=1 from unknown colorID, got %v", result["Default"])
	}
	if stored["Sage"] != 2 || stored["Default"] != 1 {
		BT.Fatalf("expected cached spendings to match result, got %+v", stored)
	}
}

func TestGetSpendingCacheMissGoogleErrorDoesNotStore(BT *testing.T) {
	BT.Parallel()

	now := time.Date(2026, 3, 5, 12, 0, 0, 0, time.UTC)
	setCalled := false
	cacheRepo := &fakeSpendingCacheRepo{
		getFn: func(context.Context, string, time.Time) (calendar.Spendings, error) {
			return nil, nil
		},
		setFn: func(context.Context, string, time.Time, calendar.Spendings) error {
			setCalled = true
			return nil
		},
	}

	svc := NewCalendarService(&fakeRepo{
		getFn: func(context.Context, string) (*calendar.CalendarConnection, error) {
			return &calendar.CalendarConnection{
				UID:          "uid",
				CalendarID:   "primary",
				AccessToken:  "token",
				RefreshToken: "refresh",
				Expiry:       now.Add(time.Hour),
			}, nil
		},
	}, cacheRepo, &fakeCalendarClient{
		listEventsFn: func(context.Context, string, string, time.Time, time.Time) ([]calendar.Event, error) {
			return nil, errors.New("google unavailable")
		},
	}, &fakeOAuth{}, "secret")

	_, err := svc.GetSpending(context.Background(), "uid", now, now.AddDate(0, 0, 6))
	if err == nil {
		BT.Fatal("expected error")
	}
	if setCalled {
		BT.Fatal("expected cache set not to be called when google fetch fails")
	}
}
