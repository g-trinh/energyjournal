package energy

import "context"

type EnergyRepository interface {
	GetByDate(ctx context.Context, uid, date string) (*EnergyLevels, error)
	Upsert(ctx context.Context, levels EnergyLevels) error
}
