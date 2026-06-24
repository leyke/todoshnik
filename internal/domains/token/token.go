package token

import "todoshnik/internal/auth"

type Token struct {
	ID     int
	UserID int `json:"user_id"`
	// возможно для hash больше подошел бы типа byte
	Hash   string          `json:"hash"`
	Device auth.DeviceType `json:"Device"`
	// почему это число а не time.Time?
	ExpiresAt int64 `json:"expires_at"`
}
