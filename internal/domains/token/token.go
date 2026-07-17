package token

import "time"

type Token struct {
	ID     int
	UserID int        `json:"user_id"`
	Hash   string     `json:"hash"`
	Device DeviceType `json:"Device"`
	// почему это число а не time.Time?
	ExpiresAt time.Time `json:"expires_at"`
}
