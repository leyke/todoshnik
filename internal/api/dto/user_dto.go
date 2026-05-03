package dto

type UserSignUpRequestDto struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
	Name     string `json:"name" validate:"required,min=3"`
}

type UserSignInRequestDto struct {
	Login    string `json:"login" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
}

type AuthResponseDto struct {
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
}
