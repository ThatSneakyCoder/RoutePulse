package service

import (
	"context"
	"errors"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IdentityService struct {
	repo domain.UserRepository
	log  *zap.SugaredLogger
}

func NewIdentityService(repo domain.UserRepository, log *zap.SugaredLogger) *IdentityService {
	return &IdentityService{
		repo: repo,
		log: log,
	}
}

func (s *IdentityService) RegisterUser(
	ctx context.Context,
	user *domain.User,
) (*domain.User, error) {

	s.log.Infow("register user request received")

	// Sanity check
	if user == nil {
		s.log.Error("register user failed: user is nil")
		return nil, errors.New("user is nil")
	}

	// Normalize input
	normalizeUser(user)

	if err := validateEmail(user.Email); err != nil {
		s.log.Warnw("email validation failed",
			"email", user.Email,
			"err", err,
		)
		return nil, err
	}

	if err := validatePassword(user.Password); err != nil {
		s.log.Warnw("password validation failed",
			"err", err,
		)
		return nil, err
	}

	hashed, err := hashPassword(user.Password)
	if err != nil {
		s.log.Errorw("password hashing failed", "err", err)
		return nil, err
	}
	user.Password = hashed

	user.ID = uuid.New()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		s.log.Errorw("failed to persist user",
			"user_id", user.ID,
			"email", user.Email,
			"err", err,
		)
		return nil, err
	}

	s.log.Infow("user registered successfully",
		"user_id", created.ID,
	)

	return created, nil
}