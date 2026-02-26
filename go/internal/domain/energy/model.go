package energy

import (
	"context"
	"time"
)

type EnergyLevels struct {
	UID                string
	Date               string
	Physical           int
	Mental             int
	Emotional          int
	SleepQuality       int
	StressLevel        int
	PhysicalActivity   string
	Nutrition          string
	SocialInteractions string
	TimeOutdoors       string
	Notes              string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type EnergyService interface {
	GetByDate(ctx context.Context, uid, date string) (*EnergyLevels, error)
	GetByDateRange(ctx context.Context, uid, from, to string) ([]EnergyLevels, error)
	Save(ctx context.Context, levels EnergyLevels) error
}
