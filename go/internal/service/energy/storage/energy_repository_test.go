package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	"energyjournal/internal/domain/energy"
	pkgerror "energyjournal/internal/pkg/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestEnergyRepository_GetByDate_ExistingDocument(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 2, 20, 9, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 2, 21, 9, 0, 0, 0, time.UTC)
	client := newFakeFirestoreClient()
	client.docs["energy_levels/uid-1_2026-02-21"] = map[string]any{
		"uid":       "uid-1",
		"date":      "2026-02-21",
		"physical":  7,
		"mental":    5,
		"emotional": 8,
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}
	repo := newEnergyRepositoryWithClient(client, time.Now)

	got, err := repo.GetByDate(context.Background(), "uid-1", "2026-02-21")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got.UID != "uid-1" || got.Date != "2026-02-21" {
		t.Fatalf("unexpected id/date: %+v", got)
	}
	if got.Physical != 7 || got.Mental != 5 || got.Emotional != 8 {
		t.Fatalf("unexpected levels: %+v", got)
	}
	if !got.CreatedAt.Equal(createdAt) || !got.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected timestamps: %+v", got)
	}
}

func TestEnergyRepository_GetByDate_MissingDocumentReturnsNotFound(t *testing.T) {
	t.Parallel()

	repo := newEnergyRepositoryWithClient(newFakeFirestoreClient(), time.Now)

	_, err := repo.GetByDate(context.Background(), "uid-1", "2026-02-21")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var notFoundErr *pkgerror.NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("expected NotFoundError, got %T", err)
	}
}

func TestEnergyRepository_Upsert_NewDocument(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 21, 11, 0, 0, 0, time.UTC)
	repo := newEnergyRepositoryWithClient(newFakeFirestoreClient(), func() time.Time { return now })

	err := repo.Upsert(context.Background(), energyFixture("uid-1", "2026-02-21", 6, 7, 8))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	doc := repo.client.(*fakeFirestoreClient).docs["energy_levels/uid-1_2026-02-21"]
	if doc == nil {
		t.Fatal("expected stored document")
	}
	if !doc["createdAt"].(time.Time).Equal(now) {
		t.Fatalf("expected createdAt %s, got %v", now, doc["createdAt"])
	}
	if !doc["updatedAt"].(time.Time).Equal(now) {
		t.Fatalf("expected updatedAt %s, got %v", now, doc["updatedAt"])
	}
}

func TestEnergyRepository_Upsert_ExistingDocumentPreservesCreatedAt(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 2, 20, 10, 0, 0, 0, time.UTC)
	firstUpdate := time.Date(2026, 2, 21, 8, 0, 0, 0, time.UTC)
	secondUpdate := time.Date(2026, 2, 21, 12, 0, 0, 0, time.UTC)
	client := newFakeFirestoreClient()
	client.docs["energy_levels/uid-1_2026-02-21"] = map[string]any{
		"uid":       "uid-1",
		"date":      "2026-02-21",
		"physical":  4,
		"mental":    5,
		"emotional": 6,
		"createdAt": createdAt,
		"updatedAt": firstUpdate,
	}

	repo := newEnergyRepositoryWithClient(client, func() time.Time { return secondUpdate })

	err := repo.Upsert(context.Background(), energyFixture("uid-1", "2026-02-21", 9, 9, 9))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	doc := client.docs["energy_levels/uid-1_2026-02-21"]
	if !doc["createdAt"].(time.Time).Equal(createdAt) {
		t.Fatalf("expected createdAt to remain %s, got %v", createdAt, doc["createdAt"])
	}
	if !doc["updatedAt"].(time.Time).Equal(secondUpdate) {
		t.Fatalf("expected updatedAt %s, got %v", secondUpdate, doc["updatedAt"])
	}
}

func energyFixture(uid, date string, physical, mental, emotional int) energy.EnergyLevels {
	return energy.EnergyLevels{
		UID:       uid,
		Date:      date,
		Physical:  physical,
		Mental:    mental,
		Emotional: emotional,
	}
}

type fakeFirestoreClient struct {
	docs map[string]map[string]any
}

func newFakeFirestoreClient() *fakeFirestoreClient {
	return &fakeFirestoreClient{docs: map[string]map[string]any{}}
}

func (f *fakeFirestoreClient) Collection(path string) collectionRef {
	return &fakeCollectionRef{client: f, collection: path}
}

type fakeCollectionRef struct {
	client     *fakeFirestoreClient
	collection string
}

func (f *fakeCollectionRef) Doc(path string) documentRef {
	return &fakeDocumentRef{
		client: f.client,
		key:    f.collection + "/" + path,
	}
}

type fakeDocumentRef struct {
	client *fakeFirestoreClient
	key    string
}

func (f *fakeDocumentRef) Get(_ context.Context) (documentSnapshot, error) {
	doc, ok := f.client.docs[f.key]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	return &fakeDocumentSnapshot{data: cloneDoc(doc)}, nil
}

func (f *fakeDocumentRef) Set(_ context.Context, data any) error {
	typed, ok := data.(map[string]any)
	if !ok {
		return errors.New("invalid data type")
	}
	f.client.docs[f.key] = cloneDoc(typed)
	return nil
}

type fakeDocumentSnapshot struct {
	data map[string]any
}

func (f *fakeDocumentSnapshot) Data() map[string]any {
	return cloneDoc(f.data)
}

func cloneDoc(doc map[string]any) map[string]any {
	clone := make(map[string]any, len(doc))
	for key, value := range doc {
		clone[key] = value
	}
	return clone
}
