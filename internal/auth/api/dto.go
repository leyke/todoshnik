package api

type UserSignUpRequestDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserSignInRequestDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponseDto struct {
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type TgLoginRequestDto struct {
	TgUserID int64  `json:"tg_user_id"`
	Name     string `json:"name"`
}
