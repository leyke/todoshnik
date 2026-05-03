package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	return token, nil
}

func HashToken(token string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func ComparePassword(hash string, input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}

	return true
}
