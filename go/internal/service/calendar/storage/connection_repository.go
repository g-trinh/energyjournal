package storage

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"energyjournal/internal/domain/calendar"
)

const connectionCollection = "calendar_connections"

type ConnectionRepository struct {
	client *firestore.Client
}

func NewConnectionRepository(client *firestore.Client) *ConnectionRepository {
	return &ConnectionRepository{client: client}
}

func (r *ConnectionRepository) Get(ctx context.Context, uid string) (*calendar.CalendarConnection, error) {
	doc, err := r.client.Collection(connectionCollection).Doc(uid).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	data := doc.Data()
	return &calendar.CalendarConnection{
		UID:          uid,
		CalendarID:   getString(data, "calendar_id"),
		AccessToken:  getString(data, "access_token"),
		RefreshToken: getString(data, "refresh_token"),
		Expiry:       getTime(data, "expiry"),
	}, nil
}

func (r *ConnectionRepository) Upsert(ctx context.Context, conn calendar.CalendarConnection) error {
	_, err := r.client.Collection(connectionCollection).Doc(conn.UID).Set(ctx, map[string]any{
		"access_token":  conn.AccessToken,
		"refresh_token": conn.RefreshToken,
		"calendar_id":   conn.CalendarID,
		"expiry":        conn.Expiry,
	})
	return err
}

func getString(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}

func getTime(data map[string]any, key string) time.Time {
	v := data[key]
	switch t := v.(type) {
	case time.Time:
		return t
	default:
		return time.Time{}
	}
}
