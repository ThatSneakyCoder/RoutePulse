package domain

import "time"

type Vehicle struct {
	ID             string
	OrganizationID string

	PlateNumber string
	VehicleType string
	Capacity    int32

	Status string

	CreatedAt time.Time
}

type Driver struct {
	ID             string
	OrganizationID string

	FirstName string
	LastName  string

	VehicleID *string

	Status string

	CreatedAt time.Time
}

type Trip struct {
	ID             string
	OrganizationID string
	VehicleID      string
	DriverID       string
	Status         string
	StartLatitude  float64
	StartLongitude float64
	EndLatitude    float64
	EndLongitude   float64
	StartTime      *time.Time
	EndTime        *time.Time
	CreatedAt      time.Time
}

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type RoutePreview struct {
	DistanceMeters  float64
	DurationSeconds float64
	Geometry        []Coordinate
}
