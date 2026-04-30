// Package encrypt provides simple symmetric encryption for sensitive
// environment values stored in envchain layers.
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the provided key is not a valid AES key length.
var ErrInvalidKey = errors.New("encrypt: key must be 16, 24, or 32 bytes")

// ErrInvalidCiphertext is returned when decryption input is malformed.
var ErrInvalidCiphertext = errors.New("encrypt: invalid ciphertext")

// Encryptor encrypts and decrypts string values using AES-GCM.
type Encryptor struct {
	key []byte
}

// NewEncryptor creates an Encryptor using the given raw key.
// key must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func NewEncryptor(key []byte) (*Encryptor, error) {
	switch len(key) {
	case 16, 24, 32:
		// valid
	default:
		return nil, ErrInvalidKey
	}
	cp := make([]byte, len(key))
	copy(cp, key)
	return &Encryptor{key: cp}, nil
}

// Encrypt encrypts plaintext and returns a base64-encoded ciphertext string.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decodes and decrypts a base64-encoded ciphertext produced by Encrypt.
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", ErrInvalidCiphertext
	}
	plain, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plain), nil
}
