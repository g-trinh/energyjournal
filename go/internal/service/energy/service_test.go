package energy

import (
	"context"
	"errors"
	"testing"
	"time"

	"energyjournal/internal/domain/energy"
	pkgerror "energyjournal/internal/pkg/error"
)

func TestService_GetByDate_ValidDateDelegatesToRepo(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{
		getByDate: func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
			return &energy.EnergyLevels{
				UID:       uid,
				Date:      date,
				Physical:  7,
				Mental:    6,
				Emotional: 8,
			}, nil
		},
	}
	svc := NewEnergyService(repo)

	got, err := svc.GetByDate(context.Background(), "uid-1", "2026-02-21")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got == nil || got.Physical != 7 {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestService_GetByDate_InvalidDateReturnsValidationError(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{}
	svc := NewEnergyService(repo)

	_, err := svc.GetByDate(context.Background(), "uid-1", "2026/02/21")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var validationErr *pkgerror.InputValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected InputValidationError, got %T", err)
	}
}

func TestService_Save_ValidInputSetsUpdatedAtAndDelegates(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 21, 12, 0, 0, 0, time.UTC)
	repo := &mockEnergyRepository{}
	svc := newServiceWithClock(repo, func() time.Time { return now })

	createdAt := time.Date(2026, 2, 20, 9, 0, 0, 0, time.UTC)
	err := svc.Save(context.Background(), energy.EnergyLevels{
		UID:       "uid-1",
		Date:      "2026-02-21",
		Physical:  7,
		Mental:    5,
		Emotional: 8,
		CreatedAt: createdAt,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if repo.lastSaved == nil {
		t.Fatal("expected repository Upsert to be called")
	}
	if !repo.lastSaved.UpdatedAt.Equal(now) {
		t.Fatalf("expected UpdatedAt %s, got %s", now, repo.lastSaved.UpdatedAt)
	}
	if !repo.lastSaved.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected CreatedAt %s, got %s", createdAt, repo.lastSaved.CreatedAt)
	}
}

func TestService_Save_DimensionOutOfRangeReturnsValidationError(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{}
	svc := NewEnergyService(repo)

	err := svc.Save(context.Background(), energy.EnergyLevels{
		UID:       "uid-1",
		Date:      "2026-02-21",
		Physical:  -1,
		Mental:    5,
		Emotional: 6,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var validationErr *pkgerror.InputValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected InputValidationError, got %T", err)
	}
}

func TestService_Save_MalformedDateReturnsValidationError(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{}
	svc := NewEnergyService(repo)

	err := svc.Save(context.Background(), energy.EnergyLevels{
		UID:       "uid-1",
		Date:      "invalid-date",
		Physical:  4,
		Mental:    5,
		Emotional: 6,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var validationErr *pkgerror.InputValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected InputValidationError, got %T", err)
	}
}

func TestService_Save_PropagatesRepositoryErrors(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repository failure")
	repo := &mockEnergyRepository{
		upsert: func(ctx context.Context, levels energy.EnergyLevels) error {
			return repoErr
		},
	}
	svc := NewEnergyService(repo)

	err := svc.Save(context.Background(), energy.EnergyLevels{
		UID:       "uid-1",
		Date:      "2026-02-21",
		Physical:  7,
		Mental:    5,
		Emotional: 8,
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}
}

type mockEnergyRepository struct {
	getByDate func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error)
	upsert    func(ctx context.Context, levels energy.EnergyLevels) error
	lastSaved *energy.EnergyLevels
}

func (m *mockEnergyRepository) GetByDate(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
	if m.getByDate != nil {
		return m.getByDate(ctx, uid, date)
	}
	return nil, nil
}

func (m *mockEnergyRepository) Upsert(ctx context.Context, levels energy.EnergyLevels) error {
	copyLevels := levels
	m.lastSaved = &copyLevels
	if m.upsert != nil {
		return m.upsert(ctx, levels)
	}
	return nil
}
