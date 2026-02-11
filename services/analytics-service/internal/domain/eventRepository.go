package domain

import (
	"context"
	"time"
)

type EventRepository interface {
	// Read side
	CountVehiclesInTransit(ctx context.Context) (uint64, error)
	CountTripsToday(ctx context.Context) (uint64, error)

	// Write side (event ingestion)
	InsertIdentityUserRegistered(ctx context.Context, eventTime time.Time, userID string, email string) error
}
