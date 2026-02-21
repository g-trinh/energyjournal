package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"energyjournal/internal/domain/energy"
	pkgerror "energyjournal/internal/pkg/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const energyLevelsCollection = "energy_levels"

type FirestoreEnergyRepository struct {
	client  firestoreClient
	timeNow func() time.Time
}

func NewEnergyRepository(client *firestore.Client) *FirestoreEnergyRepository {
	return &FirestoreEnergyRepository{
		client:  &firestoreClientAdapter{client: client},
		timeNow: time.Now,
	}
}

func newEnergyRepositoryWithClient(client firestoreClient, timeNow func() time.Time) *FirestoreEnergyRepository {
	return &FirestoreEnergyRepository{
		client:  client,
		timeNow: timeNow,
	}
}

func (r *FirestoreEnergyRepository) GetByDate(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
	docID := energyLevelDocID(uid, date)
	snapshot, err := r.client.Collection(energyLevelsCollection).Doc(docID).Get(ctx)
	if err != nil {
		if isNotFound(err) {
			return nil, pkgerror.NewNotFoundError("energy_levels", docID)
		}
		return nil, err
	}

	return docToEnergyLevels(snapshot), nil
}

func (r *FirestoreEnergyRepository) Upsert(ctx context.Context, levels energy.EnergyLevels) error {
	docID := energyLevelDocID(levels.UID, levels.Date)
	docRef := r.client.Collection(energyLevelsCollection).Doc(docID)

	createdAt := time.Time{}
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		if !isNotFound(err) {
			return err
		}
		createdAt = r.timeNow()
	} else {
		createdAt, _ = getTimestamp(snapshot.Data(), "createdAt")
		if createdAt.IsZero() {
			createdAt = r.timeNow()
		}
	}

	updatedAt := r.timeNow()
	payload := map[string]any{
		"uid":       levels.UID,
		"date":      levels.Date,
		"physical":  levels.Physical,
		"mental":    levels.Mental,
		"emotional": levels.Emotional,
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}

	return docRef.Set(ctx, payload)
}

func energyLevelDocID(uid, date string) string {
	return fmt.Sprintf("%s_%s", uid, date)
}

func docToEnergyLevels(snapshot documentSnapshot) *energy.EnergyLevels {
	data := snapshot.Data()

	createdAt, _ := getTimestamp(data, "createdAt")
	updatedAt, _ := getTimestamp(data, "updatedAt")

	return &energy.EnergyLevels{
		UID:       getString(data, "uid"),
		Date:      getString(data, "date"),
		Physical:  getInt(data, "physical"),
		Mental:    getInt(data, "mental"),
		Emotional: getInt(data, "emotional"),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func isNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}

func getString(data map[string]any, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

func getInt(data map[string]any, key string) int {
	value, ok := data[key]
	if !ok || value == nil {
		return 0
	}

	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	default:
		return 0
	}
}

func getTimestamp(data map[string]any, key string) (time.Time, error) {
	value, ok := data[key]
	if !ok || value == nil {
		return time.Time{}, nil
	}

	switch typed := value.(type) {
	case time.Time:
		return typed, nil
	default:
		return time.Time{}, nil
	}
}

type firestoreClient interface {
	Collection(path string) collectionRef
}

type collectionRef interface {
	Doc(path string) documentRef
}

type documentRef interface {
	Get(ctx context.Context) (documentSnapshot, error)
	Set(ctx context.Context, data any) error
}

type documentSnapshot interface {
	Data() map[string]any
}

type firestoreClientAdapter struct {
	client *firestore.Client
}

func (a *firestoreClientAdapter) Collection(path string) collectionRef {
	return &collectionRefAdapter{ref: a.client.Collection(path)}
}

type collectionRefAdapter struct {
	ref *firestore.CollectionRef
}

func (a *collectionRefAdapter) Doc(path string) documentRef {
	return &documentRefAdapter{ref: a.ref.Doc(path)}
}

type documentRefAdapter struct {
	ref *firestore.DocumentRef
}

func (a *documentRefAdapter) Get(ctx context.Context) (documentSnapshot, error) {
	snapshot, err := a.ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &documentSnapshotAdapter{snapshot: snapshot}, nil
}

func (a *documentRefAdapter) Set(ctx context.Context, data any) error {
	_, err := a.ref.Set(ctx, data)
	return err
}

type documentSnapshotAdapter struct {
	snapshot *firestore.DocumentSnapshot
}

func (a *documentSnapshotAdapter) Data() map[string]any {
	return a.snapshot.Data()
}
