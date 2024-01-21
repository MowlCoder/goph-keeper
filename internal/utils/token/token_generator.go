package token

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

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
