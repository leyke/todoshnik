package accesstoken

import "todoshnik/internal/domain"

type AccessTokenRepositoryInterface interface {
	GetAllByUserID(id int) []*domain.Token
	GetUserIDByToken(token string) int
	GetExpiredTokens() []*domain.Token

	Create(token *domain.Token) (*domain.Token, error)
	Delete(token *domain.Token) error
}
