package repository

import (
	"context"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
)

type UserStore struct {
	// db *sql.DB (later)
}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (s *UserStore) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
 return nil, nil
}
