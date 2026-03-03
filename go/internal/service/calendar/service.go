package calendar

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"energyjournal/internal/domain/calendar"
	integgoogle "energyjournal/internal/integration/google"
	errpkg "energyjournal/internal/pkg/error"
)

var colorNamesByID = map[string]string{
	"":   "Default",
	"1":  "Tomato",
	"2":  "Flamingo",
	"3":  "Tangerine",
	"4":  "Banana",
	"5":  "Sage",
	"6":  "Basil",
	"7":  "Peacock",
	"8":  "Blueberry",
	"9":  "Lavender",
	"10": "Grape",
	"11": "Graphite",
}

type oauthProvider interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}

type googleCalendarClient interface {
	ListCalendars(ctx context.Context, token string) ([]calendar.CalendarItem, error)
	ListEvents(ctx context.Context, token, calendarID string, start, end time.Time) ([]integgoogle.Event, error)
}

type CalendarService struct {
	repo         calendar.CalendarConnectionRepository
	googleClient googleCalendarClient
	oauth        oauthProvider
	stateSecret  string
	stateTTL     time.Duration
	now          func() time.Time
}

func NewCalendarService(repo calendar.CalendarConnectionRepository, googleClient googleCalendarClient, oauth oauthProvider, stateSecret string) *CalendarService {
	return &CalendarService{
		repo:         repo,
		googleClient: googleClient,
		oauth:        oauth,
		stateSecret:  stateSecret,
		stateTTL:     15 * time.Minute,
		now:          time.Now,
	}
}

func (s *CalendarService) GetStatus(ctx context.Context, uid string) (calendar.ConnectionStatus, error) {
	conn, err := s.repo.Get(ctx, uid)
	if err != nil {
		return "", err
	}
	if conn == nil {
		return calendar.StatusDisconnected, nil
	}
	if conn.CalendarID == "" {
		return calendar.StatusPendingSelection, nil
	}
	return calendar.StatusConnected, nil
}

func (s *CalendarService) BuildAuthURL(uid string) string {
	state := s.signState(uid, s.now())
	return s.oauth.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
}

func (s *CalendarService) HandleCallback(ctx context.Context, code, state string) error {
	uid, err := s.verifyState(state)
	if err != nil {
		return errpkg.NewInputValidationError("state", "invalid state")
	}

	token, err := s.oauth.Exchange(ctx, code)
	if err != nil {
		return errpkg.NewInputValidationError("code", "invalid authorization code")
	}

	return s.repo.Upsert(ctx, calendar.CalendarConnection{
		UID:          uid,
		CalendarID:   "",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	})
}

func (s *CalendarService) GetCalendars(ctx context.Context, uid string) ([]calendar.CalendarItem, error) {
	conn, err := s.requireConnection(ctx, uid)
	if err != nil {
		return nil, err
	}
	if conn.AccessToken == "" {
		return nil, errpkg.NewCalendarNotConnectedError("google calendar not connected")
	}
	return s.googleClient.ListCalendars(ctx, conn.AccessToken)
}

func (s *CalendarService) SetCalendar(ctx context.Context, uid, calendarID string) error {
	conn, err := s.requireConnection(ctx, uid)
	if err != nil {
		return err
	}
	if conn.AccessToken == "" {
		return errpkg.NewCalendarNotConnectedError("google calendar not connected")
	}

	conn.CalendarID = calendarID
	return s.repo.Upsert(ctx, *conn)
}

func (s *CalendarService) GetSpending(ctx context.Context, uid string, start, end time.Time) (calendar.Spendings, error) {
	conn, err := s.requireConnection(ctx, uid)
	if err != nil {
		return nil, err
	}
	if conn.AccessToken == "" || conn.CalendarID == "" {
		return nil, errpkg.NewCalendarNotConnectedError("google calendar not connected")
	}

	token := &oauth2.Token{
		AccessToken:  conn.AccessToken,
		RefreshToken: conn.RefreshToken,
		Expiry:       conn.Expiry,
	}

	if !conn.Expiry.IsZero() && !s.now().Before(conn.Expiry) {
		refreshed, refreshErr := s.oauth.TokenSource(ctx, token).Token()
		if refreshErr != nil {
			return nil, refreshErr
		}
		conn.AccessToken = refreshed.AccessToken
		if refreshed.RefreshToken != "" {
			conn.RefreshToken = refreshed.RefreshToken
		}
		conn.Expiry = refreshed.Expiry
		if err := s.repo.Upsert(ctx, *conn); err != nil {
			return nil, err
		}
		token = refreshed
	}

	events, err := s.googleClient.ListEvents(ctx, token.AccessToken, conn.CalendarID, start, end)
	if err != nil {
		return nil, err
	}

	out := calendar.Spendings{}
	for _, event := range events {
		if event.Start.IsZero() || event.End.IsZero() || !event.End.After(event.Start) {
			continue
		}
		name, ok := colorNamesByID[event.ColorID]
		if !ok {
			name = "Default"
		}
		out[name] += event.End.Sub(event.Start).Hours()
	}

	return out, nil
}

func (s *CalendarService) requireConnection(ctx context.Context, uid string) (*calendar.CalendarConnection, error) {
	conn, err := s.repo.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errpkg.NewCalendarNotConnectedError("google calendar not connected")
	}
	return conn, nil
}

func (s *CalendarService) signState(uid string, now time.Time) string {
	payload := uid + "|" + strconv.FormatInt(now.Unix(), 10)
	mac := hmac.New(sha256.New, []byte(s.stateSecret))
	_, _ = mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))
	return base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + signature))
}

func (s *CalendarService) verifyState(state string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return "", errors.New("invalid parts")
	}
	uid := parts[0]
	ts := parts[1]
	sig := parts[2]

	payload := uid + "|" + ts
	mac := hmac.New(sha256.New, []byte(s.stateSecret))
	_, _ = mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return "", errors.New("invalid signature")
	}

	unixSeconds, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return "", err
	}
	issuedAt := time.Unix(unixSeconds, 0)
	if s.now().Sub(issuedAt) > s.stateTTL {
		return "", fmt.Errorf("state expired")
	}

	return uid, nil
}
