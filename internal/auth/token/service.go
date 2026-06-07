package token

import (
	"context"
	"os"
	"strconv"
	"time"

	"todoshnik/internal/auth"
	apperrors "todoshnik/internal/errors"
	"todoshnik/internal/user"
)

const (
	DefaultTokenTtl int = 14 // количество дней жизни токена
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Add(ctx context.Context, user *user.User, device auth.DeviceType) (string, error) {
	token, tokenError := GenerateToken()
	if tokenError != nil {
		return "", tokenError
	}

	// 1. SALT и TOKEN_TTL_DAYS должны быть настройками компонента
	// 2. Разве соль не должна быть уникальной для пользователя? В целом не силен
	hash := HashToken(token, os.Getenv("SALT"))

	// тут лучше использовать тип Duration и местод time.ParseDuration()
	tokenTtl, tokenError := strconv.Atoi(os.Getenv("TOKEN_TTL_DAYS"))
	// дефолты тоже нужно определять не в бизнес логике, а на этапе инициализации конфига, эта логика тут не нужна
	// да и с точки зрения производительности каждый раз читать читать переменные окружения, делать парсинг это ненужный
	// оверхед
	if tokenError != nil {
		tokenTtl = DefaultTokenTtl
	}
	localTime := time.Now().AddDate(0, 0, tokenTtl).Unix()

	newToken := &Token{
		UserID:    user.ID,
		Hash:      hash,
		ExpiresAt: localTime,
		Device:    device,
	}

	newToken, err := s.repo.Create(ctx, newToken)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) Get(ctx context.Context, rawToken string) (*Token, error) {
	hash := HashToken(rawToken, os.Getenv("SALT"))

	token, err := s.repo.GetByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, apperrors.ErrNotFound
	}

	return token, nil
}

func (s *Service) ClearExpiredTokens(ctx context.Context) int {
	tokens, err := s.repo.GetExpiredTokens(ctx)
	if err != nil {
		return 0
	}
	counter := 0
	for _, token := range tokens {
		err := s.repo.Delete(ctx, token)
		// потеря ошибки, можно вернуть массив ошибок или сделать join
		if err != nil {
			counter++
		}
	}

	return counter
}
