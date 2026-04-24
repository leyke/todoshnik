package domain

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name" validate:"required,min=2"`
	TelegramID int64  `json:"telegram_id"`
}