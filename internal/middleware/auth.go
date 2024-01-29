package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/utils/usercontext"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
)

type tokenParser interface {
	Parse(token string) (*domain.TokenClaim, error)
}

// AuthMiddleware - struct responsible for validation user session
type AuthMiddleware struct {
	tokenParser tokenParser
}

// NewAuthMiddleware - constructor for AuthMiddleware struct
func NewAuthMiddleware(tokenParser tokenParser) *AuthMiddleware {
	return &AuthMiddleware{
		tokenParser: tokenParser,
	}
}

// Middleware - responsible for checking request headers for validating JWT token
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromHeader(r.Header)
		if err != nil {
			httputils.SendStatusCode(w, http.StatusUnauthorized)
			return
		}

		claim, err := m.tokenParser.Parse(token)
		if err != nil {
			httputils.SendStatusCode(w, http.StatusUnauthorized)
			return
		}

		ctx := usercontext.SetUserIDToContext(r.Context(), claim.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTokenFromHeader(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")

	if authorization == "" {
		return "", fmt.Errorf("header not found")
	}

	parts := strings.Split(authorization, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid token")
	}

	return parts[1], nil
}
