package password

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct{}

func NewHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (h *PasswordHasher) Hash(original string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(original), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *PasswordHasher) Equal(original string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(original)) == nil
}
