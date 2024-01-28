package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHash(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		err  error
	}{
		{
			name: "valid",
			val:  "password",
			err:  nil,
		},
	}

	for _, testCase := range testCases {
		hasher := NewHasher()

		t.Run(testCase.name, func(t *testing.T) {
			hash, err := hasher.Hash(testCase.val)
			assert.Equal(t, testCase.err, err)
			if err == nil {
				assert.NotEqual(t, 0, len(hash))
			}
		})
	}
}

func TestPasswordEqual(t *testing.T) {
	testCases := []struct {
		name            string
		passwordToHash  string
		passwordToEqual string
		result          bool
	}{
		{
			name:            "valid",
			passwordToHash:  "password",
			passwordToEqual: "password",
			result:          true,
		},
		{
			name:            "not valid",
			passwordToHash:  "password123",
			passwordToEqual: "password",
			result:          false,
		},
	}

	for _, testCase := range testCases {
		hasher := NewHasher()

		t.Run(testCase.name, func(t *testing.T) {
			hash, err := hasher.Hash(testCase.passwordToHash)
			require.NoError(t, err)

			isEqual := hasher.Equal(testCase.passwordToEqual, hash)
			assert.Equal(t, testCase.result, isEqual)
		})
	}
}
