package domain

import "github.com/golang-jwt/jwt/v4"

type TokenClaim struct {
	jwt.RegisteredClaims
	ID int `json:"id"`
}
