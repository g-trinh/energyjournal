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
	client  *firestore.Client
	timeNow func() time.Time
}

func NewEnergyRepository(client *firestore.Client) *FirestoreEnergyRepository {
	return &FirestoreEnergyRepository{
		client:  client,
		timeNow: time.Now,
	}
}

func (r *FirestoreEnergyRepository) GetByDate(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
	docID := energyLevelDocID(uid, date)
	snapshot, err := r.client.Collection(energyLevelsCollection).Doc(docID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, pkgerror.NewNotFoundError("energy_levels", docID)
		}
		return nil, err
	}

	data := snapshot.Data()
	return &energy.EnergyLevels{
		UID:       getString(data, "uid"),
		Date:      getString(data, "date"),
		Physical:  getInt(data, "physical"),
		Mental:    getInt(data, "mental"),
		Emotional: getInt(data, "emotional"),
		CreatedAt: getTimestamp(data, "createdAt"),
		UpdatedAt: getTimestamp(data, "updatedAt"),
	}, nil
}

func (r *FirestoreEnergyRepository) Upsert(ctx context.Context, levels energy.EnergyLevels) error {
	docID := energyLevelDocID(levels.UID, levels.Date)
	docRef := r.client.Collection(energyLevelsCollection).Doc(docID)

	createdAt := r.timeNow()
	snapshot, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return err
		}
	} else {
		if existing := getTimestamp(snapshot.Data(), "createdAt"); !existing.IsZero() {
			createdAt = existing
		}
	}

	_, err = docRef.Set(ctx, map[string]any{
		"uid":       levels.UID,
		"date":      levels.Date,
		"physical":  levels.Physical,
		"mental":    levels.Mental,
		"emotional": levels.Emotional,
		"createdAt": createdAt,
		"updatedAt": r.timeNow(),
	})
	return err
}

func energyLevelDocID(uid, date string) string {
	return fmt.Sprintf("%s_%s", uid, date)
}

func getString(data map[string]any, key string) string {
	v, _ := data[key].(string)
	return v
}

func getInt(data map[string]any, key string) int {
	switch v := data[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func getTimestamp(data map[string]any, key string) time.Time {
	v, _ := data[key].(time.Time)
	return v
}
