package user

import (
	"context"
	"errors"

	appuser "todoshnik/internal/domains/user"
	usererrors "todoshnik/internal/domains/user/errors"

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

func (repo *DBRepository) List(ctx context.Context) ([]*appuser.User, error) {
	var result []*appuser.User
	err := repo.db.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *DBRepository) GetByID(ctx context.Context, id int) (*appuser.User, error) {
	var item *appuser.User
	result := repo.db.WithContext(ctx).First(&item, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, usererrors.ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *DBRepository) GetByLogin(ctx context.Context, login string) (*appuser.User, error) {
	var item *appuser.User
	result := repo.db.WithContext(ctx).First(&item, "login = ?", login)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, usererrors.ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *DBRepository) GetByTgId(ctx context.Context, userTgId int64) (*appuser.User, error) {
	var item *appuser.User
	result := repo.db.WithContext(ctx).First(&item, "telegram_id = ?", userTgId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, usererrors.ErrNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (repo *DBRepository) Create(ctx context.Context, user *appuser.User) (*appuser.User, error) {
	result := repo.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *DBRepository) Update(ctx context.Context, user *appuser.User) error {
	result := repo.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DBRepository) Delete(ctx context.Context, user *appuser.User) error {
	result := repo.db.WithContext(ctx).Delete(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
