package user

import (
	"context"

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

func (repo *DBRepository) List(ctx context.Context) []*User {
	var result []*User
	repo.db.WithContext(ctx).Find(&result)

	return result
}

func (repo *DBRepository) GetByID(ctx context.Context, id int) (*User, bool) {
	var item *User
	result := repo.db.WithContext(ctx).First(&item, id)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *DBRepository) GetByLogin(ctx context.Context, login string) (*User, bool) {
	var item *User
	result := repo.db.WithContext(ctx).First(&item, "login = ?", login)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *DBRepository) GetByTgId(ctx context.Context, userTgId int64) (*User, bool) {
	var item *User
	result := repo.db.WithContext(ctx).First(&item, "telegram_id = ?", userTgId)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *DBRepository) Create(ctx context.Context, user *User) (*User, error) {
	result := repo.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *DBRepository) Update(ctx context.Context, user *User) error {
	result := repo.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DBRepository) Delete(ctx context.Context, user *User) error {
	result := repo.db.WithContext(ctx).Delete(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
