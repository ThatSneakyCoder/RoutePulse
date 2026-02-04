package service

import (
	"context"
	"errors"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityService struct {
	repo domain.UserRepository
	log  *zap.SugaredLogger
	jwt  *auth.JWTAuthenticator
}

func NewIdentityService(repo domain.UserRepository, log *zap.SugaredLogger, jwt *auth.JWTAuthenticator) *IdentityService {
	return &IdentityService{
		repo: repo,
		log:  log,
		jwt:  jwt,
	}
}

func (s *IdentityService) RegisterUser(
	ctx context.Context,
	user *domain.User,
	plainPassword string,
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

	if err := user.Password.Set(plainPassword); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

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

func (s *IdentityService) Login(
	ctx context.Context,
	email string,
	plainPassword string,
) (*domain.User, *auth.Token, error) {

	s.log.Debugw("login attempt started", "email", email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Errorw("login failed: user lookup error",
			"email", email,
			"err", err,
		)
		return nil, nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if err := user.Password.Compare(plainPassword); err != nil {
		s.log.Errorw("login failed: password mismatch",
			"user_id", user.ID,
		)
		return nil, nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	now := time.Now()
	expiresAt := now.Add(s.jwt.Expiry())

	claims := jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": expiresAt.Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"iss": s.jwt.Issuer(),
		"aud": s.jwt.Audience(),
	}

	tokenStr, err := s.jwt.GenerateToken(claims)
	if err != nil {
		s.log.Errorw("login failed: jwt generation error",
			"user_id", user.ID,
			"err", err,
		)
		return nil, nil, status.Error(codes.Internal, "failed to generate token")
	}

	s.log.Infow("login successful",
		"user_id", user.ID,
		"expires_at", expiresAt,
	)

	return user, &auth.Token{
		AccessToken: tokenStr,
		ExpiresAt:   expiresAt,
	}, nil
}
