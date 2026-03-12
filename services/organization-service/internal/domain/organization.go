package domain

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID
	Name        string
	OwnerUserID uuid.UUID
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrganizationMember struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	JoinedAt       time.Time
}
