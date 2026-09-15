package user

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name" validate:"required,min=2"`
	TelegramID   int64  `json:"telegram_id"`
	Login        string `json:"login" validate:"required,min=2"`
	PasswordHash string `json:"password_hash"`
}
