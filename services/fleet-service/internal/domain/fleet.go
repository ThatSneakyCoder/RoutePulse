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
	StartTime      *time.Time
	EndTime        *time.Time
	CreatedAt      time.Time
}