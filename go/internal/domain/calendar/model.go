package calendar

import "time"

// Spendings represents time spent per event type (in hours).
type Spendings map[string]float64

// SpendingService defines the contract for spending retrieval.
type SpendingService interface {
	GetSpending(start, end time.Time) (Spendings, error)
}
