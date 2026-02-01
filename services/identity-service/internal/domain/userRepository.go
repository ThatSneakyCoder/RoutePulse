package domain

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
}
