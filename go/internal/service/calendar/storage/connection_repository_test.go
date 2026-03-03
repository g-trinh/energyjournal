package storage

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"energyjournal/internal/domain/calendar"
)

type fakeStore struct {
	docs map[string]map[string]any
}

func newFakeStore() *fakeStore {
	return &fakeStore{docs: map[string]map[string]any{}}
}

func (s *fakeStore) Doc(uid string) docStore {
	return &fakeDoc{uid: uid, store: s}
}

type fakeDoc struct {
	uid   string
	store *fakeStore
}

func (d *fakeDoc) Get(_ context.Context) (map[string]any, error) {
	data, ok := d.store.docs[d.uid]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing")
	}
	return data, nil
}

func (d *fakeDoc) Set(_ context.Context, data map[string]any) error {
	d.store.docs[d.uid] = data
	return nil
}

func TestConnectionRepository_GetReturnsNilWhenMissing(t *testing.T) {
	t.Parallel()

	repo := &ConnectionRepository{store: newFakeStore()}

	conn, err := repo.Get(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if conn != nil {
		t.Fatalf("expected nil connection, got %+v", conn)
	}
}

func TestConnectionRepository_UpsertThenGet(t *testing.T) {
	t.Parallel()

	repo := &ConnectionRepository{store: newFakeStore()}
	expectedExpiry := time.Date(2026, time.March, 3, 12, 0, 0, 0, time.UTC)

	err := repo.Upsert(context.Background(), calendar.CalendarConnection{
		UID:          "uid-2",
		CalendarID:   "primary",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Expiry:       expectedExpiry,
	})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	conn, err := repo.Get(context.Background(), "uid-2")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if conn == nil {
		t.Fatal("expected non-nil connection")
	}
	if conn.UID != "uid-2" {
		t.Fatalf("expected uid uid-2, got %s", conn.UID)
	}
	if conn.CalendarID != "primary" {
		t.Fatalf("expected calendar_id primary, got %s", conn.CalendarID)
	}
	if conn.AccessToken != "access-token" {
		t.Fatalf("expected access token access-token, got %s", conn.AccessToken)
	}
	if conn.RefreshToken != "refresh-token" {
		t.Fatalf("expected refresh token refresh-token, got %s", conn.RefreshToken)
	}
	if !conn.Expiry.Equal(expectedExpiry) {
		t.Fatalf("expected expiry %s, got %s", expectedExpiry, conn.Expiry)
	}
}
