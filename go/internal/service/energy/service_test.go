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

func TestService_GetByDateRange_MissingFromDefaultsToLastSevenDays(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 23, 12, 0, 0, 0, time.UTC)
	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if from != "2026-02-17" || to != "2026-02-23" {
				t.Fatalf("unexpected default range: %s to %s", from, to)
			}
			return []energy.EnergyLevels{}, nil
		},
	}
	svc := newServiceWithClock(repo, func() time.Time { return now })

	_, err := svc.GetByDateRange(context.Background(), "uid-1", "", "2026-02-21")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestService_GetByDateRange_InvalidDateDefaultsToLastSevenDays(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 23, 12, 0, 0, 0, time.UTC)
	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if from != "2026-02-17" || to != "2026-02-23" {
				t.Fatalf("unexpected default range: %s to %s", from, to)
			}
			return []energy.EnergyLevels{}, nil
		},
	}
	svc := newServiceWithClock(repo, func() time.Time { return now })

	_, err := svc.GetByDateRange(context.Background(), "uid-1", "2026-02-22", "bad-date")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestService_GetByDateRange_FromAfterToDefaultsToLastSevenDays(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 2, 23, 12, 0, 0, 0, time.UTC)
	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if from != "2026-02-17" || to != "2026-02-23" {
				t.Fatalf("unexpected default range: %s to %s", from, to)
			}
			return []energy.EnergyLevels{}, nil
		},
	}
	svc := newServiceWithClock(repo, func() time.Time { return now })

	_, err := svc.GetByDateRange(context.Background(), "uid-1", "2026-02-22", "2026-02-21")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestService_GetByDateRange_Valid_Delegates(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if from != "2026-02-01" || to != "2026-02-14" {
				t.Fatalf("unexpected range: %s to %s", from, to)
			}
			return []energy.EnergyLevels{{UID: uid, Date: from, Physical: 6, Mental: 5, Emotional: 7}}, nil
		},
	}
	svc := NewEnergyService(repo)

	got, err := svc.GetByDateRange(context.Background(), "uid-1", "2026-02-01", "2026-02-14")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestService_GetByDateRange_ExactlyThirtyDays(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if to != "2026-01-31" {
				t.Fatalf("expected to remain 2026-01-31, got %s", to)
			}
			return []energy.EnergyLevels{}, nil
		},
	}
	svc := NewEnergyService(repo)

	_, err := svc.GetByDateRange(context.Background(), "uid-1", "2026-01-01", "2026-01-31")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestService_GetByDateRange_ThirtyOneDays_Clamps(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			if to != "2026-01-31" {
				t.Fatalf("expected clamp to 2026-01-31, got %s", to)
			}
			return []energy.EnergyLevels{}, nil
		},
	}
	svc := NewEnergyService(repo)

	_, err := svc.GetByDateRange(context.Background(), "uid-1", "2026-01-01", "2026-02-01")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestService_GetByDateRange_RepoError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repository failure")
	repo := &mockEnergyRepository{
		getByDateRange: func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
			return nil, repoErr
		},
	}
	svc := NewEnergyService(repo)

	_, err := svc.GetByDateRange(context.Background(), "uid-1", "2026-02-01", "2026-02-14")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
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

func TestService_Save_ContextEnumsValidPass(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{}
	svc := NewEnergyService(repo)

	err := svc.Save(context.Background(), energy.EnergyLevels{
		UID:                "uid-1",
		Date:               "2026-02-21",
		Physical:           4,
		Mental:             5,
		Emotional:          6,
		SleepQuality:       intPtr(5),
		StressLevel:        intPtr(1),
		PhysicalActivity:   "moderate",
		Nutrition:          "good",
		SocialInteractions: "positive",
		TimeOutdoors:       "30min_1hr",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestService_Save_InvalidContextEnumReturnsValidationError(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{}
	svc := NewEnergyService(repo)

	err := svc.Save(context.Background(), energy.EnergyLevels{
		UID:              "uid-1",
		Date:             "2026-02-21",
		Physical:         4,
		Mental:           5,
		Emotional:        6,
		PhysicalActivity: "sprint",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var validationErr *pkgerror.InputValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected InputValidationError, got %T", err)
	}
}

func TestService_Save_SleepQualityValidation(t *testing.T) {
	t.Parallel()

	repo := &mockEnergyRepository{}
	svc := NewEnergyService(repo)

	base := energy.EnergyLevels{
		UID:       "uid-1",
		Date:      "2026-02-21",
		Physical:  4,
		Mental:    5,
		Emotional: 6,
	}

	err := svc.Save(context.Background(), base)
	if err != nil {
		t.Fatalf("expected nil sleepQuality to pass, got %v", err)
	}

	withFive := base
	withFive.SleepQuality = intPtr(5)
	err = svc.Save(context.Background(), withFive)
	if err != nil {
		t.Fatalf("expected sleepQuality=5 to pass, got %v", err)
	}

	withSix := base
	withSix.SleepQuality = intPtr(6)
	err = svc.Save(context.Background(), withSix)
	if err == nil {
		t.Fatal("expected error for sleepQuality=6, got nil")
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

func intPtr(n int) *int { return &n }

type mockEnergyRepository struct {
	getByDate      func(ctx context.Context, uid, date string) (*energy.EnergyLevels, error)
	getByDateRange func(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error)
	upsert         func(ctx context.Context, levels energy.EnergyLevels) error
	lastSaved      *energy.EnergyLevels
}

func (m *mockEnergyRepository) GetByDate(ctx context.Context, uid, date string) (*energy.EnergyLevels, error) {
	if m.getByDate != nil {
		return m.getByDate(ctx, uid, date)
	}
	return nil, nil
}

func (m *mockEnergyRepository) GetByDateRange(ctx context.Context, uid, from, to string) ([]energy.EnergyLevels, error) {
	if m.getByDateRange != nil {
		return m.getByDateRange(ctx, uid, from, to)
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
