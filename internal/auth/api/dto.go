package api

type UserSignUpRequestDto struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
}

type UserSignInRequestDto struct {
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

type AuthResponseDto struct {
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token,omitempty"`
}

type TgLoginRequestDto struct {
	TgUserID int64  `json:"tg_user_id"`
	Name     string `json:"name"`
}
