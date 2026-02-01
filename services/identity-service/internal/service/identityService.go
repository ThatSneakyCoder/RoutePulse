package service

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
)

type IdentityService struct {
	repo domain.UserRepository
}

func NewIdentityService(repo domain.UserRepository) *IdentityService {
	return &IdentityService{
		repo: repo,
	}
}

func (s *IdentityService) RegisterUser(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {
	// business rules here
	return s.repo.Create(ctx, user)
}
