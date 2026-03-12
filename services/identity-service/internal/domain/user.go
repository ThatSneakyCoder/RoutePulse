package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`

	Password password `json:"-"`

	IsActive   bool `json:"is_active"`
	IsVerified bool `json:"is_verified"`

	EmailVerification *EmailVerification `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmailVerification struct {
	TokenHash  *string
	ExpiresAt  *time.Time
}

