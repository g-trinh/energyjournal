package energy

import "context"

type EnergyRepository interface {
	GetByDate(ctx context.Context, uid, date string) (*EnergyLevels, error)
	GetByDateRange(ctx context.Context, uid, from, to string) ([]EnergyLevels, error)
	Upsert(ctx context.Context, levels EnergyLevels) error
}
