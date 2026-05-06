package accesstoken

import (
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

func (repo *AccessTokenDbRepository) GetAllByUserID(userID int) []*domain.Token {
	var result []*domain.Token
	repo.db.Where("user_id = ?", userID).Find(&result)
	return result
}

func (repo *AccessTokenDbRepository) GetUserIDByToken(hash string) int {
	var token *domain.Token
	repo.db.Where("hash = ?", hash).First(&token)
	return token.UserID
}

func (repo *AccessTokenDbRepository) GetExpiredTokens() []*domain.Token {
	var result []*domain.Token
	localTime := time.Now().Unix()
	repo.db.Where("expires_at < ?", localTime).Find(&result)
	return result
}

func (repo *AccessTokenDbRepository) Create(token *domain.Token) (*domain.Token, error) {
	result := repo.db.Create(token)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (repo *AccessTokenDbRepository) Delete(token *domain.Token) error {
	result := repo.db.Delete(token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
