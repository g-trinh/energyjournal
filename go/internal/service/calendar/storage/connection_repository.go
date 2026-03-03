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

type docStore interface {
	Get(ctx context.Context) (map[string]any, error)
	Set(ctx context.Context, data map[string]any) error
}

type connectionStore interface {
	Doc(uid string) docStore
}

type firestoreConnectionStore struct {
	client *firestore.Client
}

func (s *firestoreConnectionStore) Doc(uid string) docStore {
	return &firestoreDocStore{ref: s.client.Collection(connectionCollection).Doc(uid)}
}

type firestoreDocStore struct {
	ref *firestore.DocumentRef
}

func (s *firestoreDocStore) Get(ctx context.Context) (map[string]any, error) {
	doc, err := s.ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	return doc.Data(), nil
}

func (s *firestoreDocStore) Set(ctx context.Context, data map[string]any) error {
	_, err := s.ref.Set(ctx, data)
	return err
}

type ConnectionRepository struct {
	store connectionStore
}

func NewConnectionRepository(client *firestore.Client) *ConnectionRepository {
	return &ConnectionRepository{
		store: &firestoreConnectionStore{client: client},
	}
}

func (r *ConnectionRepository) Get(ctx context.Context, uid string) (*calendar.CalendarConnection, error) {
	data, err := r.store.Doc(uid).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	return &calendar.CalendarConnection{
		UID:          uid,
		CalendarID:   getString(data, "calendar_id"),
		AccessToken:  getString(data, "access_token"),
		RefreshToken: getString(data, "refresh_token"),
		Expiry:       getTime(data, "expiry"),
	}, nil
}

func (r *ConnectionRepository) Upsert(ctx context.Context, conn calendar.CalendarConnection) error {
	return r.store.Doc(conn.UID).Set(ctx, map[string]any{
		"access_token":  conn.AccessToken,
		"refresh_token": conn.RefreshToken,
		"calendar_id":   conn.CalendarID,
		"expiry":        conn.Expiry,
	})
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
