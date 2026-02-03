package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooWeak = errors.New("password does not meet requirements")

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return ErrPasswordTooWeak
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
