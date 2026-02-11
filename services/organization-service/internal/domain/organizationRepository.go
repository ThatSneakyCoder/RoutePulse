package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) (*Organization, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Organization, error)
}
