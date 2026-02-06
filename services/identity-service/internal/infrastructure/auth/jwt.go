package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
	exp    time.Duration
}

func (a *JWTAuthenticator) Issuer() string {
	return a.iss
}

func (a *JWTAuthenticator) Audience() string {
	return a.aud
}

func (a *JWTAuthenticator) Expiry() time.Duration {
	return a.exp
}

func NewJWTAuthenticator(secret, aud, iss string, exp time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud, exp}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
