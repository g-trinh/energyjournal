package energy

import (
	"context"
	"regexp"
	"time"

	domain "energyjournal/internal/domain/energy"
	pkgerror "energyjournal/internal/pkg/error"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

var physicalActivityValues = map[string]struct{}{
	"none":     {},
	"light":    {},
	"moderate": {},
	"intense":  {},
}

var nutritionValues = map[string]struct{}{
	"poor":      {},
	"average":   {},
	"good":      {},
	"excellent": {},
}

var socialInteractionValues = map[string]struct{}{
	"negative": {},
	"neutral":  {},
	"positive": {},
}

var timeOutdoorsValues = map[string]struct{}{
	"none":        {},
	"under_30min": {},
	"30min_1hr":   {},
	"over_1hr":    {},
}

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
	fromDate, fromOK := parseDate(from)
	toDate, toOK := parseDate(to)
	if !fromOK || !toOK || toDate.Before(fromDate) {
		toDate = s.timeNow()
		fromDate = toDate.AddDate(0, 0, -6)
	}

	from = fromDate.Format("2006-01-02")
	to = toDate.Format("2006-01-02")

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
	if err := validateOptionalScaleField("sleepQuality", levels.SleepQuality); err != nil {
		return err
	}
	if err := validateOptionalScaleField("stressLevel", levels.StressLevel); err != nil {
		return err
	}
	if err := validateOptionalEnumField("physicalActivity", levels.PhysicalActivity, physicalActivityValues); err != nil {
		return err
	}
	if err := validateOptionalEnumField("nutrition", levels.Nutrition, nutritionValues); err != nil {
		return err
	}
	if err := validateOptionalEnumField("socialInteractions", levels.SocialInteractions, socialInteractionValues); err != nil {
		return err
	}
	if err := validateOptionalEnumField("timeOutdoors", levels.TimeOutdoors, timeOutdoorsValues); err != nil {
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

func parseDate(date string) (time.Time, bool) {
	if date == "" || !datePattern.MatchString(date) {
		return time.Time{}, false
	}
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

func validateLevel(field string, value int) error {
	if value < 0 || value > 10 {
		return pkgerror.NewInputValidationError(field, "must be between 0 and 10")
	}
	return nil
}

func validateOptionalScaleField(field string, value *int) error {
	if value == nil {
		return nil
	}
	if *value < 1 || *value > 5 {
		return pkgerror.NewInputValidationError(field, "must be between 1 and 5")
	}
	return nil
}

func validateOptionalEnumField(field, value string, allowed map[string]struct{}) error {
	if value == "" {
		return nil
	}
	if _, ok := allowed[value]; !ok {
		return pkgerror.NewInputValidationError(field, "invalid value")
	}
	return nil
}
