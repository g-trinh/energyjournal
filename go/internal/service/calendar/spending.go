package calendar

import (
	"time"

	"energyjournal/internal/domain/calendar"
)

type spendingService struct{}

// NewSpendingService creates a new SpendingService implementation.
func NewSpendingService() calendar.SpendingService {
	return &spendingService{}
}

// GetSpending returns hardcoded dummy spendings.
func (s *spendingService) GetSpending(start, end time.Time) (calendar.Spendings, error) {
	return calendar.Spendings{
		"Travail": 35,
		"Perso":   5.6,
		"Routine": 7.75,
		"Repas":   10,
	}, nil
}
