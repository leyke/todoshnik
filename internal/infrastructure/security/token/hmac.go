package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type HMACHasher struct {
	secret string
}

func NewHMACHasher(secret string) *HMACHasher {
	return &HMACHasher{
		secret: secret,
	}
}

func (h *HMACHasher) Generate() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (h *HMACHasher) Hash(token string) string {
	hash := hmac.New(sha256.New, []byte(h.secret))

	hash.Write([]byte(token))

	return hex.EncodeToString(hash.Sum(nil))
}
