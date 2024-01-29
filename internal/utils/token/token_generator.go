package token

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

// Generator - struct responsible for generate JWT tokens
type Generator struct {
}

// NewGenerator - constructor for Generator struct
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate - generate JWT token from domain.User struct
func (g *Generator) Generate(ctx context.Context, user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.TokenClaim{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24)),
		},
	})

	tokenStr, err := token.SignedString(getTokenSecretKey())
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
