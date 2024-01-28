package token

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokenSecret(t *testing.T) {
	testCases := []struct {
		name    string
		prepare func()
		result  []byte
	}{
		{
			name: "get from env",
			prepare: func() {
				os.Unsetenv(jwtSecretEnvKey)
				os.Setenv(jwtSecretEnvKey, "SECRET TOKEN")
			},
			result: []byte("SECRET TOKEN"),
		},
		{
			name: "get default value",
			prepare: func() {
				os.Unsetenv(jwtSecretEnvKey)
			},
			result: []byte(defaultSecretVal),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.prepare()
			secret := getTokenSecretKey()
			assert.Equal(t, testCase.result, secret)
		})
	}
}
