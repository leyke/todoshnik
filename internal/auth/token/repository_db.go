package token

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDbRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{
		db: db,
	}
}

func (repo *DBRepository) GetAllByUserID(ctx context.Context, userID int) []*Token {
	var result []*Token
	repo.db.WithContext(ctx).Where("user_id = ?", userID).Find(&result)
	return result
}

func (repo *DBRepository) GetByHash(ctx context.Context, hash string) *Token {
	var token *Token
	repo.db.WithContext(ctx).Where("hash = ?", hash).First(&token)
	return token
}

func (repo *DBRepository) GetExpiredTokens(ctx context.Context) []*Token {
	var result []*Token
	localTime := time.Now().Unix()
	repo.db.WithContext(ctx).Where("expires_at < ?", localTime).Find(&result)
	return result
}

func (repo *DBRepository) Create(ctx context.Context, token *Token) (*Token, error) {
	result := repo.db.WithContext(ctx).Create(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (repo *DBRepository) Delete(ctx context.Context, token *Token) error {
	result := repo.db.WithContext(ctx).Delete(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
