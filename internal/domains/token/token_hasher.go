package token

type TokenGenerator interface {
	Generate() (string, error)
	Hash(token string) string
}
