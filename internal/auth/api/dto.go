package api

//  1. изучи какие есть стандартные аннотации у маршалера, и обрати внимание на omitempty
//  2. думаю пользователи не обрадуются если ты будет логировать куда-то в файл или stdout их пароли
//     как минимум на маскировать/солить+хешировать/закрывать звездочками, но в чистом виде не в коем случаем не писать
//     это хуже чем креды в коде
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
