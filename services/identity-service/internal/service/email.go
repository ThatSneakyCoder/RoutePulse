package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
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

// func sendVerifyTokenEmailToUser(email string) error {

// }

func generateVerifyEmailToken() (string, error) {
	const digits = "0123456789"
	const length = 6

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		b[i] = digits[n.Int64()]
	}

	return string(b), nil
}

func hashVerifyEmailToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
