package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooWeak = errors.New("password does not meet requirements")

type password struct {
	text *string
	hash []byte
}

func (p *password) SetHash(hash []byte) {
	p.hash = hash
	p.text = nil
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return ErrPasswordTooWeak
	}

	return nil
}

func (p *password) Set(plain string) error {
	if err := validatePassword(plain); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.SetHash(hash)

	return nil
}

func (p *password) Compare(plain string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
}
