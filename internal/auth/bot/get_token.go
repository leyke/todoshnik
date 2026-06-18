package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	tokenerror "todoshnik/internal/auth/token/errors"
)

const tgAuthLoginURL string = "/auth/tg/login"

func (h *Handler) GetToken(ctx context.Context, tgUser TgLoginRequestDto) (*AuthInfo, error) {

	// TODO расставить таймауты
	// 3. Распространенная ошибка– в сетевых клиентах не устанавливать таймаут. В продакшн коде должен быть:
	// - прокинут контекст с дедлайном или таймаутом
	// - настроены таймауты через изменяемый без перезагрузки конфиг, или по крайней мере изменяемый при перезагрузке
	// - возможно настроены ретраи, возможно с экспоненциальным бекофом
	//   (может быть вот это https://pkg.go.dev/github.com/cenkalti/backoff/v4)
	// - может быть какие-то ручки деградации/circuit breaker
	response, err := h.api.Post(
		ctx,
		tgAuthLoginURL,
		tgUser,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected status code: %d",
			response.StatusCode,
		)
	}

	var responseInfo UserAuthInfoResponseDto
	if err := json.NewDecoder(response.Body).Decode(&responseInfo); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if responseInfo.AccessToken == "" {
		return nil, fmt.Errorf(
			"%w for user id %d",
			tokenerror.ErrNotFound,
			responseInfo.UserID,
		)
	}

	return &AuthInfo{
		UserID:      responseInfo.UserID,
		AccessToken: responseInfo.AccessToken,
	}, nil
}
