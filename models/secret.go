package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"time"
	"io"
)

type Secret struct {
	ID           string    `json:"id"`
	EncryptedMsg []byte    `json:"encrypted_msg"`
	IV           []byte    `json:"iv"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	OneTime      bool      `json:"one_time"`
	CreatedBy    string    `json:"created_by"`
	PasswordHash string    `json:"password_hash"`
}

// GenerateID returns a random 12-character string
func GenerateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 12)
	for i := range b {
		rand.Read(b[i : i+1])
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// Encrypt plaintext using AES
func Encrypt(plaintext []byte) ([]byte, []byte, error) {
	key := []byte("verysecretkey123") // 🔐 Replace with env-based key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	return ciphertext, iv, nil
}

// Decrypt ciphertext using AES
func Decrypt(ciphertext []byte, iv []byte) ([]byte, error) {
	key := []byte("verysecretkey123") // 🔐 Replace with env-based key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
