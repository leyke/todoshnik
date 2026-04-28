package domain

type Token struct {
	ID        int
	UserID    int
	Hash      string
	ExpiresAt int64
}
