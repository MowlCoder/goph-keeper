package token

import (
	"github.com/golang-jwt/jwt/v4"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tokenString string) (*domain.TokenClaim, error) {
	tokenClaim := &domain.TokenClaim{}
	token, err := jwt.ParseWithClaims(tokenString, tokenClaim, func(token *jwt.Token) (interface{}, error) {
		return getTokenSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return tokenClaim, nil
}
