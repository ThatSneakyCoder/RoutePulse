package service

import (
	"context"
	"errors"
	"time"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/auth"
	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/infrastructure/mailer"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityService struct {
	repo   domain.UserRepository
	log    *zap.SugaredLogger
	jwt    *auth.JWTAuthenticator
	mailer mailer.Client
}

func NewIdentityService(repo domain.UserRepository, log *zap.SugaredLogger, jwt *auth.JWTAuthenticator, mailer mailer.Client) *IdentityService {
	return &IdentityService{
		repo:   repo,
		log:    log,
		jwt:    jwt,
		mailer: mailer,
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

	verifyEmailToken, err := generateVerifyEmailToken()
	if err != nil {
		s.log.Warnw("verifyEmailToken generation failed",
			"email", user.Email,
			"err", err,
		)
		return nil, err
	}

	hashedVerifyEmailToken := hashVerifyEmailToken(verifyEmailToken)
	const verifyEmailTokenTTL = 10 * time.Minute
	expiresAt := time.Now().UTC().Add(verifyEmailTokenTTL)

	user.EmailVerification = &domain.EmailVerification{
		TokenHash: &hashedVerifyEmailToken,
		ExpiresAt: &expiresAt,
	}

	// user is active by default
	user.IsActive = true

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

	// 3. Send email to user with token
	data := struct {
		Email     string
		Code      string
		ExpiresIn int // minutes
		Year      int
	}{
		Email:     user.Email,
		Code:      verifyEmailToken,
		ExpiresIn: 10,
		Year:      time.Now().Year(),
	}

	if err := s.mailer.Send(
		"status_email.tmpl",
		user.Email,
		data,
	); err != nil {
		s.log.Errorw("failed to send email", "err", err)
	}

	s.log.Infow("status code mail sent to user successfully",
		"data.Email", data.Email,
	)

	return created, nil
}

func (s *IdentityService) VerifyUserEmail(ctx context.Context, email, token string) (string, bool, error) {
	s.log.Infow("verify user email requested",
		"email", email,
	)

	hashedToken := hashVerifyEmailToken(token)
	now := time.Now().UTC()

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Errorw("email verification failed: user lookup error",
			"email", email,
			"err", err,
		)
		return "", false, err
	}

	if user.IsVerified {
		// idempotent behavior
		return user.ID.String(), true, nil
	}

	if user.EmailVerification == nil {
		return "", false, status.Error(codes.PermissionDenied, "no verification pending")
	}

	if user.EmailVerification.ExpiresAt.Before(now) {
		// delete user from database
		s.repo.DeleteUserByEmail(ctx, user.Email)
		return "", false, status.Error(codes.PermissionDenied, "verification token expired")
	}

	if *user.EmailVerification.TokenHash != hashedToken {
		return "", false, status.Error(codes.PermissionDenied, "invalid verification token")
	}

	// mark verified
	user.IsVerified = true
	user.EmailVerification = nil
	user.UpdatedAt = now

	if err := s.repo.MarkEmailVerified(ctx, user); err != nil {
		return "", false, err
	}

	s.log.Infow("user email verified successfully",
		"user_id", user.ID,
		"email", email,
	)

	return user.ID.String(), true, nil
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

func (s *IdentityService) ValidateJWTTokenAndGetUser(
	ctx context.Context,
	accessToken string,
) (*domain.User, error) {

	s.log.Debugw("validating jwt access token")

	jwtToken, err := s.jwt.ValidateToken(accessToken)
	if err != nil {
		s.log.Errorw("jwt validation failed", "err", err)
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		s.log.Errorw("invalid jwt claims type")
		return nil, status.Error(codes.Unauthenticated, "invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		s.log.Errorw("missing sub claim in token")
		return nil, status.Error(codes.Unauthenticated, "invalid subject in token")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		s.log.Errorw("invalid subject format",
			"sub", sub,
			"err", err,
		)
		return nil, status.Error(codes.Unauthenticated, "malformed subject")
	}

	s.log.Debugw("fetching user for token",
		"user_id", userID,
	)

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		s.log.Errorw("user lookup failed",
			"user_id", userID,
			"err", err,
		)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	s.log.Debugw("user resolved from token",
		"user_id", user.ID,
	)

	return user, nil
}

func (s *IdentityService) ForgotPassword(
	ctx context.Context,
	email string,
) error {

	s.log.Infow("forgot password requested", "email", email)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// prevent user enumeration
		return nil
	}

	token, err := generatePasswordResetToken()
	if err != nil {
		return err
	}

	hashed := hashPasswordResetToken(token)
	expiresAt := time.Now().UTC().Add(10 * time.Minute)

	reset := &domain.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: hashed,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.Upsert(ctx, reset); err != nil {
		return err
	}

	data := struct {
		Email     string
		Code      string
		ExpiresIn int
		Year      int
	}{
		Email:     user.Email,
		Code:      token,
		ExpiresIn: 10,
		Year:      time.Now().Year(),
	}

	if err := s.mailer.Send(
		"password_reset.tmpl",
		user.Email,
		data,
	); err != nil {
		s.log.Errorw("failed to send reset password email", "err", err)
	}

	s.log.Infow("forgot password mail sent successfully", "email", email)

	return nil
}

func (s *IdentityService) ResetPassword(
	ctx context.Context,
	email, token, newPassword string,
) error {
	s.log.Infow("reset password requested", "email", email, "token", token)

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return status.Error(codes.PermissionDenied, "invalid reset request")
	}

	reset, err := s.repo.GetByUserID(ctx, user.ID)
	if err != nil {
		return status.Error(codes.PermissionDenied, "no reset requested")
	}

	// If reset token has expired 10 minute mark
	if reset.ExpiresAt.Before(time.Now().UTC()) {
		_ = s.repo.DeleteByUserID(ctx, user.ID)
		return status.Error(codes.PermissionDenied, "reset token expired")
	}

	// verify whether the token user sent in their request is the token that was the same token sent to them
	hashed := hashPasswordResetToken(token)
	if reset.TokenHash != hashed {
		return status.Error(codes.PermissionDenied, "invalid reset token")
	}

	// set new password
	if err := user.Password.Set(newPassword); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.repo.UpdatePassword(ctx, user); err != nil {
		return err
	}

	_ = s.repo.DeleteByUserID(ctx, user.ID)

	s.log.Infow("password reset was successful",
		"user", user.Email,
	)

	return nil
}
