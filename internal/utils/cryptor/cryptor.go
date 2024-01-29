package cryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Cryptor - structure responsible for encryption and decryption strings or bytes arrays
type Cryptor struct {
	encryptionKey []byte
}

// New - constructor for Cryptor structure
func New(key string) *Cryptor {
	return &Cryptor{
		encryptionKey: []byte(key),
	}
}

// Encrypt - encrypt given string and return encrypted string
func (c *Cryptor) Encrypt(raw string) (string, error) {
	b, err := c.EncryptBytes([]byte(raw))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Decrypt - decrypt given string and return decrypted string
func (c *Cryptor) Decrypt(crypted string) (string, error) {
	b, err := c.DecryptBytes([]byte(crypted))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// EncryptBytes - encrypt given byte array and return encrypted byte array
func (c *Cryptor) EncryptBytes(raw []byte) ([]byte, error) {
	rawBytes := raw
	generatedCipher, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, errors.New("error generating cipher")
	}

	gcm, err := cipher.NewGCM(generatedCipher)
	if err != nil {
		return nil, errors.New("error generating GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("error generating secured-sequence")
	}

	encryptedPassword := gcm.Seal(nonce, nonce, rawBytes, nil)
	b64Password := base64.StdEncoding.EncodeToString(encryptedPassword)

	return []byte(b64Password), nil
}

// DecryptBytes - decrypt given byte array and return decrypted byte array
func (c *Cryptor) DecryptBytes(crypted []byte) ([]byte, error) {
	encryptedPassword, err := base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, errors.New("error decoding base64 encrypted password string")
	}

	generatedCipher, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		return nil, errors.New("error generating cipher")
	}

	gcm, err := cipher.NewGCM(generatedCipher)
	if err != nil {
		return nil, errors.New("error generating GCM")
	}

	nonce, ciphertext := encryptedPassword[:gcm.NonceSize()], encryptedPassword[gcm.NonceSize():]

	decryptedText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("error attempting to decrypt AES-encrypted password")
	}

	return decryptedText, nil
}
