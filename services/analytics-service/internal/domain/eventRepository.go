package domain

import (
	"context"
)

type EventRepository interface {
	// Read side
	CountVehiclesInTransit(ctx context.Context) (uint64, error)
	CountTripsToday(ctx context.Context) (uint64, error)

	// Write side (event ingestion)
	InsertEvent(ctx context.Context, event *AnalyticsEvent) error
}
