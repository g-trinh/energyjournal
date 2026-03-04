package google

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestListCalendars(t *testing.T) {
	t.Parallel()

	client := &GoogleCalendarClient{
		baseURL: "https://example.test/calendar/v3",
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.Path != "/calendar/v3/users/me/calendarList" {
					t.Fatalf("unexpected path: %s", req.URL.Path)
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"items": [
							{"id":"primary","summary":"Main","backgroundColor":"#8fa58b"}
						]
					}`)),
					Header: make(http.Header),
				}, nil
			}),
		},
	}

	calendars, err := client.ListCalendars(context.Background(), "token")
	if err != nil {
		t.Fatalf("ListCalendars returned error: %v", err)
	}
	if len(calendars) != 1 {
		t.Fatalf("expected 1 calendar, got %d", len(calendars))
	}
	if calendars[0].ID != "primary" || calendars[0].Name != "Main" || calendars[0].Color != "#8fa58b" {
		t.Fatalf("unexpected calendar item: %+v", calendars[0])
	}
}

func TestListEvents(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 31, 23, 59, 59, 0, time.UTC)

	client := &GoogleCalendarClient{
		baseURL: "https://example.test/calendar/v3",
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if req.URL.Path != "/calendar/v3/calendars/primary/events" {
					t.Fatalf("unexpected path: %s", req.URL.Path)
				}
				if req.URL.Query().Get("singleEvents") != "true" {
					t.Fatalf("singleEvents not set")
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"items": [
							{
								"colorId":"5",
								"start":{"dateTime":"2026-03-03T09:00:00Z"},
								"end":{"dateTime":"2026-03-03T11:30:00Z"}
							}
						]
					}`)),
					Header: make(http.Header),
				}, nil
			}),
		},
	}

	events, err := client.ListEvents(context.Background(), "token", "primary", start, end)
	if err != nil {
		t.Fatalf("ListEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ColorID != "5" {
		t.Fatalf("expected colorId 5, got %s", events[0].ColorID)
	}
	if got := events[0].End.Sub(events[0].Start); got != 150*time.Minute {
		t.Fatalf("expected 150m duration, got %s", got)
	}
}

func TestListCalendarsNon2xxReturnsTypedError(t *testing.T) {
	t.Parallel()

	client := &GoogleCalendarClient{
		baseURL: "https://example.test/calendar/v3",
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"upstream failure"}}`)),
					Header:     make(http.Header),
				}, nil
			}),
		},
	}

	_, err := client.ListCalendars(context.Background(), "token")
	if err == nil {
		t.Fatal("expected error")
	}

	typed, ok := err.(*GoogleAPIError)
	if !ok {
		t.Fatalf("expected *GoogleAPIError, got %T", err)
	}
	if typed.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected status 502, got %d", typed.StatusCode)
	}
	if typed.Body != "upstream failure" {
		t.Fatalf("expected message upstream failure, got %q", typed.Body)
	}
}

func TestListEventsSkipsInvalidDateTime(t *testing.T) {
	t.Parallel()

	client := &GoogleCalendarClient{
		baseURL: "https://example.test/calendar/v3",
		httpClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"items": [
							{
								"colorId":"5",
								"start":{"dateTime":"invalid"},
								"end":{"dateTime":"2026-03-03T11:30:00Z"}
							}
						]
					}`)),
					Header: make(http.Header),
				}, nil
			}),
		},
	}

	events, err := client.ListEvents(context.Background(), "token", "primary", time.Now(), time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("ListEvents returned error: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no events after invalid date filtering, got %d", len(events))
	}
}

func TestGoogleAPIErrorError(t *testing.T) {
	t.Parallel()
	err := (&GoogleAPIError{StatusCode: 418}).Error()
	expected := "google calendar api returned status 418"
	if err != expected {
		t.Fatalf("expected %q, got %q", expected, err)
	}
}

func ExampleGoogleAPIError_Error() {
	fmt.Println((&GoogleAPIError{StatusCode: 500}).Error())
	// Output: google calendar api returned status 500
}
