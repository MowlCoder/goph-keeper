package token

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

func TestParser_Parse(t *testing.T) {
	testCases := []struct {
		name  string
		user  *domain.User
		isErr bool
	}{
		{
			name:  "valid",
			user:  &domain.User{ID: 1},
			isErr: false,
		},
		{
			name:  "not valid",
			user:  nil,
			isErr: true,
		},
	}

	for _, testCase := range testCases {
		generator := NewGenerator()
		parser := NewParser()

		t.Run(testCase.name, func(t *testing.T) {
			var token string
			var err error

			if testCase.user != nil {
				token, err = generator.Generate(context.Background(), *testCase.user)
				require.NoError(t, err)
			}

			tokenClaim, err := parser.Parse(token)
			if testCase.isErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.user.ID, tokenClaim.ID)
			}
		})
	}
}
