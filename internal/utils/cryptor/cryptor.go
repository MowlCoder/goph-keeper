package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Cryptor struct {
	encryptionKey []byte
}

func New(key string) *Cryptor {
	return &Cryptor{
		encryptionKey: []byte(key),
	}
}

func (c *Cryptor) Encrypt(raw string) (string, error) {
	rawBytes := []byte(raw)
	generatedCipher, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return "", errors.New("error generating cipher")
	}

	gcm, err := cipher.NewGCM(generatedCipher)
	if err != nil {
		return "", errors.New("error generating GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("error generating secured-sequence")
	}

	encryptedPassword := gcm.Seal(nonce, nonce, rawBytes, nil)
	b64Password := base64.StdEncoding.EncodeToString(encryptedPassword)

	return b64Password, nil
}

func (c *Cryptor) Decrypt(crypted string) (string, error) {
	encryptedPassword, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", errors.New("error decoding base64 encrypted password string")
	}

	generatedCipher, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return "", errors.New("error generating cipher")
	}

	gcm, err := cipher.NewGCM(generatedCipher)
	if err != nil {
		return "", errors.New("error generating GCM")
	}

	nonce, ciphertext := encryptedPassword[:gcm.NonceSize()], encryptedPassword[gcm.NonceSize():]

	decryptedText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("error attempting to decrypt AES-encrypted password")
	}

	return string(decryptedText), nil
}
