package user

type TgLoginRequestDto struct {
	TgUserID int64  `json:"tg_user_id"`
	Name     string `json:"name"`
}

type UserAuthInfoResponseDto struct {
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
}
