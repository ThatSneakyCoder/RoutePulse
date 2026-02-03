package service

import (
	"errors"
	"strings"

	"github.com/ThatSneakyCoder/RoutePulse/services/identity-service/internal/domain"
)

var ErrInvalidEmail = errors.New("invalid email address")

func normalizeUser(u *domain.User) {
	u.Email = normalizeEmail(u.Email)
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateEmail(email string) error {
	if email == "" || !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}
	return nil
}
