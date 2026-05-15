package auth

import (
	"golang.org/x/crypto/bcrypt"
)

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
