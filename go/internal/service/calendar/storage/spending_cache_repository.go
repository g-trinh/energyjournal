package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"energyjournal/internal/domain/calendar"
)

const spendingCacheCollection = "spending_cache"

type SpendingCacheRepository struct {
	client *firestore.Client
}

type spendingCacheDocument struct {
	Data      map[string]float64 `firestore:"data"`
	CreatedAt time.Time          `firestore:"created_at"`
}

func NewSpendingCacheRepository(client *firestore.Client) *SpendingCacheRepository {
	return &SpendingCacheRepository{client: client}
}

func (r *SpendingCacheRepository) Get(ctx context.Context, uid string, weekStart time.Time) (calendar.Spendings, error) {
	docID := buildSpendingCacheDocID(uid, weekStart)
	doc, err := r.client.Collection(spendingCacheCollection).Doc(docID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	var payload spendingCacheDocument
	if err := doc.DataTo(&payload); err != nil {
		return nil, err
	}
	if payload.Data == nil {
		return nil, nil
	}

	out := make(calendar.Spendings, len(payload.Data))
	for label, hours := range payload.Data {
		out[label] = hours
	}
	return out, nil
}

func (r *SpendingCacheRepository) Set(ctx context.Context, uid string, weekStart time.Time, spendings calendar.Spendings) error {
	docID := buildSpendingCacheDocID(uid, weekStart)
	payload := spendingCacheDocument{
		Data:      make(map[string]float64, len(spendings)),
		CreatedAt: time.Now().UTC(),
	}
	for label, hours := range spendings {
		payload.Data[label] = hours
	}

	_, err := r.client.Collection(spendingCacheCollection).Doc(docID).Set(ctx, payload)
	return err
}

func buildSpendingCacheDocID(uid string, weekStart time.Time) string {
	return fmt.Sprintf("%s_%s", uid, weekStart.Format("2006-01-02"))
}
