package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"

	"energyjournal/internal/domain/calendar"
)

const googleCalendarAPIBaseURL = "https://www.googleapis.com/calendar/v3"

type GoogleAPIError struct {
	StatusCode int
	Body       string
}

func (e *GoogleAPIError) Error() string {
	return fmt.Sprintf("google calendar api returned status %d", e.StatusCode)
}

type Event struct {
	ColorID string
	Start   time.Time
	End     time.Time
}

type GoogleCalendarClient struct {
	baseURL    string
	transport  http.RoundTripper
	httpClient *http.Client
}

func NewGoogleCalendarClient() *GoogleCalendarClient {
	return &GoogleCalendarClient{
		baseURL: googleCalendarAPIBaseURL,
	}
}

func (c *GoogleCalendarClient) ListCalendars(ctx context.Context, token string) ([]calendar.CalendarItem, error) {
	endpoint := fmt.Sprintf("%s/users/me/calendarList", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.oauthClient(ctx, token).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeGoogleAPIError(resp)
	}

	var payload struct {
		Items []struct {
			ID             string `json:"id"`
			Summary        string `json:"summary"`
			BackgroundColor string `json:"backgroundColor"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	items := make([]calendar.CalendarItem, 0, len(payload.Items))
	for _, item := range payload.Items {
		items = append(items, calendar.CalendarItem{
			ID:    item.ID,
			Name:  item.Summary,
			Color: item.BackgroundColor,
		})
	}

	return items, nil
}

func (c *GoogleCalendarClient) ListEvents(ctx context.Context, token, calendarID string, start, end time.Time) ([]Event, error) {
	path := fmt.Sprintf("%s/calendars/%s/events", c.baseURL, url.PathEscape(calendarID))
	query := url.Values{}
	query.Set("timeMin", start.Format(time.RFC3339))
	query.Set("timeMax", end.Format(time.RFC3339))
	query.Set("singleEvents", "true")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.oauthClient(ctx, token).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, decodeGoogleAPIError(resp)
	}

	var payload struct {
		Items []struct {
			ColorID string `json:"colorId"`
			Start   struct {
				DateTime string `json:"dateTime"`
			} `json:"start"`
			End struct {
				DateTime string `json:"dateTime"`
			} `json:"end"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	events := make([]Event, 0, len(payload.Items))
	for _, item := range payload.Items {
		startAt, err := time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			continue
		}
		endAt, err := time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			continue
		}
		events = append(events, Event{
			ColorID: item.ColorID,
			Start:   startAt,
			End:     endAt,
		})
	}

	return events, nil
}

func (c *GoogleCalendarClient) oauthClient(ctx context.Context, token string) *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}

	source := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := oauth2.NewClient(ctx, source)
	if c.transport != nil {
		client.Transport = c.transport
	}
	return client
}

func decodeGoogleAPIError(resp *http.Response) error {
	var payload struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&payload)
	return &GoogleAPIError{
		StatusCode: resp.StatusCode,
		Body:       payload.Error.Message,
	}
}
