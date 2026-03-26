package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("tracking resource not found")

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type TripCurrentLocation struct {
	TripID     string
	DriverID   string
	VehicleID  string
	Latitude   float64
	Longitude  float64
	RecordedAt time.Time
	Connection string
}

type TripLocationHistoryPoint struct {
	Latitude   float64
	Longitude  float64
	RecordedAt time.Time
}

type TripGeometry struct {
	TripID          string
	PlannedGeometry []Coordinate
	ActualGeometry  []Coordinate
}

type TrackingLocationUpdate struct {
	TripID     string
	DriverID   string
	VehicleID  string
	Latitude   float64
	Longitude  float64
	RecordedAt time.Time
	Sequence   int
}
