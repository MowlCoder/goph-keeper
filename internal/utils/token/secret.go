package token

import "os"

const (
	jwtSecretEnvKey  = "JWT_SECRET"
	defaultSecretVal = "secret"
)

func getTokenSecretKey() []byte {
	key, ok := os.LookupEnv(jwtSecretEnvKey)

	if !ok {
		return []byte(defaultSecretVal)
	}

	return []byte(key)
}
