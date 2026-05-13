package user

import (
	"context"
	"todoshnik/internal/domain"

	"gorm.io/gorm"
)

type UserDbRepository struct {
	db *gorm.DB
}

func NewUserDbRepository(db *gorm.DB) *UserDbRepository {
	return &UserDbRepository{
		db: db,
	}
}

func (repo *UserDbRepository) List(ctx context.Context) []*domain.User {
	var result []*domain.User
	repo.db.WithContext(ctx).Find(&result)

	return result
}

func (repo *UserDbRepository) GetByID(ctx context.Context, id int) (*domain.User, bool) {
	var item *domain.User
	result := repo.db.WithContext(ctx).First(&item, id)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *UserDbRepository) GetByLogin(ctx context.Context, login string) (*domain.User, bool) {
	var item *domain.User
	result := repo.db.WithContext(ctx).First(&item, "login = ?", login)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *UserDbRepository) GetUserByTgId(ctx context.Context, userTgId int64) (*domain.User, bool) {
	var item *domain.User
	result := repo.db.WithContext(ctx).First(&item, "telegram_id = ?", userTgId)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *UserDbRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := repo.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserDbRepository) Update(ctx context.Context, user *domain.User) error {
	result := repo.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *UserDbRepository) Delete(ctx context.Context, user *domain.User) error {
	result := repo.db.WithContext(ctx).Delete(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
