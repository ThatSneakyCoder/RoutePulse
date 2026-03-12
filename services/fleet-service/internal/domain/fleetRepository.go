package domain

import "context"

type FleetRepository interface {
	CreateVehicle(ctx context.Context, v *Vehicle) error
	ListVehicles(ctx context.Context, organizationID string) ([]*Vehicle, error)
	GetVehicle(ctx context.Context, vehicleID string) (*Vehicle, error)
	UpdateVehicle(ctx context.Context, vehicleID string, plateNumber string, vehicleType string, capacity int32) error
	UpdateVehicleStatus(ctx context.Context, vehicleID string, status string) error

	CreateDriver(ctx context.Context, d *Driver) error
	ListDrivers(ctx context.Context, organizationID string) ([]*Driver, error)
	UpdateDriver(ctx context.Context, driverID string, firstName string, lastName string) error

	UpdateDriverStatus(ctx context.Context, driverID string, status string) error

	CreateTrip(ctx context.Context, trip *Trip) error
	ListTrips(ctx context.Context, orgID string) ([]*Trip, error)

	// supporter methods
	HasActiveTripForVehicle(ctx context.Context, vehicleID string) (bool, error)
	HasActiveTripForDriver(ctx context.Context, driverID string) (bool, error)
	GetDriver(ctx context.Context, driverID string) (*Driver, error)
}
