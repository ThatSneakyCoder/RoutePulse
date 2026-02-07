package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	MarkEmailVerified(ctx context.Context, user *User) error
	DeleteUserByEmail(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, user *User) error
	Upsert(ctx context.Context, reset *PasswordReset) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*PasswordReset, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
