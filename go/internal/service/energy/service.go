package energy

import (
	"context"
	"regexp"
	"time"

	domain "energyjournal/internal/domain/energy"
	pkgerror "energyjournal/internal/pkg/error"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

type service struct {
	repo    domain.EnergyRepository
	timeNow func() time.Time
}

func NewEnergyService(repo domain.EnergyRepository) domain.EnergyService {
	return &service{
		repo:    repo,
		timeNow: time.Now,
	}
}

func newServiceWithClock(repo domain.EnergyRepository, timeNow func() time.Time) *service {
	return &service{
		repo:    repo,
		timeNow: timeNow,
	}
}

func (s *service) GetByDate(ctx context.Context, uid, date string) (*domain.EnergyLevels, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}

	return s.repo.GetByDate(ctx, uid, date)
}

func (s *service) GetByDateRange(ctx context.Context, uid, from, to string) ([]domain.EnergyLevels, error) {
	if err := validateDateField("from", from); err != nil {
		return nil, err
	}
	if err := validateDateField("to", to); err != nil {
		return nil, err
	}

	fromDate, _ := time.Parse("2006-01-02", from)
	toDate, _ := time.Parse("2006-01-02", to)
	if toDate.Before(fromDate) {
		return nil, pkgerror.NewInputValidationError("from", "must be on or before to")
	}

	if int(toDate.Sub(fromDate).Hours()/24) > 30 {
		toDate = fromDate.AddDate(0, 0, 30)
		to = toDate.Format("2006-01-02")
	}

	return s.repo.GetByDateRange(ctx, uid, from, to)
}

func (s *service) Save(ctx context.Context, levels domain.EnergyLevels) error {
	if err := validateDate(levels.Date); err != nil {
		return err
	}

	if err := validateLevel("physical", levels.Physical); err != nil {
		return err
	}
	if err := validateLevel("mental", levels.Mental); err != nil {
		return err
	}
	if err := validateLevel("emotional", levels.Emotional); err != nil {
		return err
	}

	levels.UpdatedAt = s.timeNow()
	return s.repo.Upsert(ctx, levels)
}

func validateDate(date string) error {
	return validateDateField("date", date)
}

func validateDateField(field, date string) error {
	if date == "" {
		return pkgerror.NewInputValidationError(field, "is required")
	}

	if !datePattern.MatchString(date) {
		return pkgerror.NewInputValidationError(field, "invalid date format, expected YYYY-MM-DD")
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return pkgerror.NewInputValidationError(field, "invalid date")
	}

	return nil
}

func validateLevel(field string, value int) error {
	if value < 0 || value > 10 {
		return pkgerror.NewInputValidationError(field, "must be between 0 and 10")
	}
	return nil
}
