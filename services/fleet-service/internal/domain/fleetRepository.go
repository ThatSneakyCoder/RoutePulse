package domain

import "context"

type FleetRepository interface {
	CreateVehicle(ctx context.Context, v *Vehicle) error
	ListVehicles(ctx context.Context, organizationID string) ([]*Vehicle, error)
}
