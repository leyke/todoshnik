package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (b *BcryptHasher) Hash(password string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *BcryptHasher) Compare(hash string, input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
