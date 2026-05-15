package token

import "todoshnik/internal/auth"

type Token struct {
	ID        int
	UserID    int             `json:"user_id"`
	Hash      string          `json:"hash"`
	Device    auth.DeviceType `json:"Device"`
	ExpiresAt int64           `json:"expires_at"`
}
