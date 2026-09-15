package user

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, input string) (bool, error)
}
