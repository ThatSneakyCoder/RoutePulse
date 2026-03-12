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
