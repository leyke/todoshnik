package accesstoken

import (
	"context"
	"time"
	"todoshnik/internal/domain"

	"gorm.io/gorm"
)

type AccessTokenDbRepository struct {
	db *gorm.DB
}

func NewAccessTokenDbRepository(db *gorm.DB) *AccessTokenDbRepository {
	return &AccessTokenDbRepository{
		db: db,
	}
}

func (repo *AccessTokenDbRepository) GetAllByUserID(ctx context.Context, userID int) []*domain.Token {
	var result []*domain.Token
	repo.db.WithContext(ctx).Where("user_id = ?", userID).Find(&result)
	return result
}

func (repo *AccessTokenDbRepository) GetUserIDByToken(ctx context.Context, hash string) int {
	var token *domain.Token
	repo.db.WithContext(ctx).Where("hash = ?", hash).First(&token)
	return token.UserID
}

func (repo *AccessTokenDbRepository) GetExpiredTokens(ctx context.Context) []*domain.Token {
	var result []*domain.Token
	localTime := time.Now().Unix()
	repo.db.WithContext(ctx).Where("expires_at < ?", localTime).Find(&result)
	return result
}

func (repo *AccessTokenDbRepository) Create(ctx context.Context, token *domain.Token) (*domain.Token, error) {
	result := repo.db.WithContext(ctx).Create(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (repo *AccessTokenDbRepository) Delete(ctx context.Context, token *domain.Token) error {
	result := repo.db.WithContext(ctx).Delete(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
