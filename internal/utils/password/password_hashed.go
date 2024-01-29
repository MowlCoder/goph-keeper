package password

import "golang.org/x/crypto/bcrypt"

// Hasher - struct responsible for hashing and comparing strings
type Hasher struct{}

// NewHasher - constructor for Hasher struct
func NewHasher() *Hasher {
	return &Hasher{}
}

// Hash - generate hash from given string
func (h *Hasher) Hash(original string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(original), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Equal - compare original string with a given hash
func (h *Hasher) Equal(original string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(original)) == nil
}
