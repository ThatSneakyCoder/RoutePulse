package service

import (
	"context"
	"errors"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/organization-service/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrganizationService struct {
	log  *zap.SugaredLogger
	repo domain.OrganizationRepository
}

func NewOrganizationService(repo domain.OrganizationRepository, log *zap.SugaredLogger) *OrganizationService {
	return &OrganizationService{
		repo: repo,
		log:  log,
	}
}

func (s *OrganizationService) CreateOrganization(
	ctx context.Context,
	name string,
	ownerUserID string,
) (*domain.Organization, error) {

	s.log.Infow("create organization request received",
		"name", name,
		"owner_user_id", ownerUserID,
	)

	if name == "" {
		s.log.Warn("organization creation failed: name is empty")
		return nil, errors.New("organization name is required")
	}

	ownerID, err := uuid.Parse(ownerUserID)
	if err != nil {
		s.log.Warnw("invalid owner user id",
			"owner_user_id", ownerUserID,
			"err", err,
		)
		return nil, err
	}

	now := time.Now().UTC()

	org := &domain.Organization{
		ID:          uuid.New(),
		Name:        name,
		OwnerUserID: ownerID,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repo.Create(ctx, org)
	if err != nil {
		s.log.Errorw("failed to persist organization",
			"organization_id", org.ID,
			"owner_user_id", ownerUserID,
			"err", err,
		)
		return nil, err
	}

	s.log.Infow("organization created successfully",
		"organization_id", created.ID,
		"owner_user_id", ownerUserID,
	)

	return created, nil
}

func (s *OrganizationService) ListOrganizationsByUserID(
	ctx context.Context,
	userID string,
) ([]*domain.Organization, error) {

	s.log.Infow("list organizations by user requested",
		"user_id", userID,
	)

	uid, err := uuid.Parse(userID)
	if err != nil {
		s.log.Warnw("invalid user id",
			"user_id", userID,
			"err", err,
		)
		return nil, err
	}

	orgs, err := s.repo.ListByUserID(ctx, uid)
	if err != nil {
		s.log.Errorw("failed to list organizations for user",
			"user_id", userID,
			"err", err,
		)
		return nil, err
	}

	s.log.Infow("organizations fetched successfully",
		"user_id", userID,
		"count", len(orgs),
	)

	return orgs, nil
}
