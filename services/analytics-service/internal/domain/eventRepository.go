package domain

import (
	"context"

)

type EventRepository interface {
	// existing
	CountVehiclesInTransit(ctx context.Context) (uint64, error)
	CountTripsToday(ctx context.Context) (uint64, error)

	// new
	CountTotalMembers(ctx context.Context) (uint64, error)
	CountActiveUsersToday(ctx context.Context) (uint64, error)
	GetRecentActivity(ctx context.Context) ([]*Event, error)

	// write
	InsertEvent(ctx context.Context, event *AnalyticsEvent) error
}
