package token

import (
	"context"
	"time"

	apptoken "todoshnik/internal/domains/token"

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

func (repo *DBRepository) GetAllByUserID(ctx context.Context, userID int) ([]*apptoken.Token, error) {
	result := make([]*apptoken.Token, 0)
	err := repo.db.WithContext(ctx).Where("user_id = ?", userID).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *DBRepository) GetByHash(ctx context.Context, hash string) (*apptoken.Token, error) {
	var token *apptoken.Token
	err := repo.db.WithContext(ctx).Where("hash = ?", hash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (repo *DBRepository) GetExpiredTokens(ctx context.Context) ([]*apptoken.Token, error) {
	var result []*apptoken.Token
	localTime := time.Now().Unix()
	err := repo.db.WithContext(ctx).Where("expires_at < ?", localTime).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *DBRepository) Create(ctx context.Context, token *apptoken.Token) (*apptoken.Token, error) {
	result := repo.db.WithContext(ctx).Create(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (repo *DBRepository) Delete(ctx context.Context, token *apptoken.Token) error {
	result := repo.db.WithContext(ctx).Delete(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
