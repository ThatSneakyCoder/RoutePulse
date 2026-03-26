package domain

import "context"

type TrackingRepository interface {
	GetTripCurrentLocation(ctx context.Context, tripID string) (*TripCurrentLocation, error)
	GetTripLocationHistory(ctx context.Context, tripID string, limit int32) ([]*TripLocationHistoryPoint, error)
	GetTripGeometry(ctx context.Context, tripID string) (*TripGeometry, error)
	StoreLocationUpdate(ctx context.Context, update TrackingLocationUpdate) error
}
