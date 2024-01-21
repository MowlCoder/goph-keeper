package token

import "os"

func getTokenSecretKey() []byte {
	key, ok := os.LookupEnv("JWT_SECRET")

	if !ok {
		return []byte("secret")
	}

	return []byte(key)
}
