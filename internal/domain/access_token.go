package domain

type DeviceType string

const (
	DeviceTypeBot DeviceType = "bot"
	DeviceTypeApi DeviceType = "api"
)

type Token struct {
	ID        int
	UserID    int        `json:"user_id"`
	Hash      string     `json:"hash"`
	Device    DeviceType `json:"Device"`
	ExpiresAt int64      `json:"expires_at"`
}
