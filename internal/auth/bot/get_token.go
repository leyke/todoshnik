package bot

import (
	"context"
	"encoding/json"
	"fmt"
)

func (h *Handler) GetToken(ctx context.Context, tgUser TgLoginRequestDto) (string, error) {
	// 1. Хардкод путей это криминал– выносим в переменные пакета/константы. Я бы рассматривал этот компонент как клиент
	// и соовтветственно разместил его на инфраструктурном слое.
	// 2. В данном случе возврат базовых скалярных типов это норм, но если клиент возвращает что-то более существенное-
	// лучше заводить доменные типы
	// плоховато что ошибка не специфицирована
	// 3. Распространенная ошибка– в сетевых клиентах не устанавливать таймаут. В продакшн коде должен быть:
	// - прокинут контекст с дедлайном или таймаутом
	// - настроены таймауты через изменяемый без перезагрузки конфиг, или по крайней мере изменяемый при перезагрузке
	// - возможно настроены ретраи, возможно с экспоненциальным бекофом
	//   (может быть вот это https://pkg.go.dev/github.com/cenkalti/backoff/v4)
	// - может быть какие-то ручки деградации/circuit breaker
	response, err := h.api.Post(
		ctx,
		"/auth/tg/login",
		tgUser,
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var responseInfo UserAuthInfoResponseDto

	err = json.NewDecoder(response.Body).Decode(&responseInfo)
	if err != nil {
		return "", err
	}

	if responseInfo.AccessToken != "" {
		return responseInfo.AccessToken, nil
	}

	return "", fmt.Errorf(
		// не надо врпапать через v, это ничего не дает, лучше %w
		"Ошибка получения токена из: %v",
		responseInfo,
	)
}
