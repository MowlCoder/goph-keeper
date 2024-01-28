package password

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(original string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(original), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *Hasher) Equal(original string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(original)) == nil
}
