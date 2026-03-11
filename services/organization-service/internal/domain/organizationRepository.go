package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *Organization) (*Organization, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Organization, error)
	ListMembersByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*OrganizationMember, error)
	GetByID(ctx context.Context, organizationID uuid.UUID) (*Organization, error)
	AddMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, role string) error
}
