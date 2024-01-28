package token

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MowlCoder/goph-keeper/internal/domain"
)

func TestGenerator_Generate(t *testing.T) {
	testCases := []struct {
		name string
		user domain.User
		err  error
	}{
		{
			name: "valid",
			user: domain.User{ID: 1},
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		generator := NewGenerator()

		t.Run(testCase.name, func(t *testing.T) {
			token, err := generator.Generate(context.Background(), testCase.user)
			assert.Equal(t, testCase.err, err)
			if err == nil {
				assert.NotEqual(t, 0, len(token))
			}
		})
	}
}
