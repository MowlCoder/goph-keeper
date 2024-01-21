package dtos

import (
	"regexp"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func validatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	hasNumber := false
	hasUpper := false
	hasLower := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsLower(c):
			hasLower = true
		default:
		}
	}

	return hasNumber && hasUpper && hasLower && hasSpecial
}

type RegisterBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *RegisterBody) Validate() bool {
	if b.Email == "" || b.Password == "" {
		return false
	}

	if !emailRegex.MatchString(b.Email) {
		return false
	}

	if !validatePassword(b.Password) {
		return false
	}

	return true
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type AuthorizeBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *AuthorizeBody) Validate() bool {
	if b.Email == "" || b.Password == "" {
		return false
	}

	return true
}

type AuthorizeResponse struct {
	Token string `json:"token"`
}
