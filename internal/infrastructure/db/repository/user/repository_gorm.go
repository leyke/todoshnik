package user

import (
	"context"
	"errors"

	appuser "todoshnik/internal/domains/user"
	usererrors "todoshnik/internal/domains/user/errors"

	"gorm.io/gorm"
)

// deprecated оставил для примера работы с gorm
type GormDBRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormDBRepository {
	return &GormDBRepository{
		db: db,
	}
}

func (repo *GormDBRepository) List(ctx context.Context) ([]*appuser.User, error) {
	var result []*appuser.User
	err := repo.db.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *GormDBRepository) GetByID(ctx context.Context, id int) (*appuser.User, error) {
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

func (repo *GormDBRepository) GetByLogin(ctx context.Context, login string) (*appuser.User, error) {
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

func (repo *GormDBRepository) GetByTgId(ctx context.Context, userTgId int64) (*appuser.User, error) {
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

func (repo *GormDBRepository) Create(ctx context.Context, user *appuser.User) (*appuser.User, error) {
	result := repo.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *GormDBRepository) Update(ctx context.Context, user *appuser.User) error {
	result := repo.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *GormDBRepository) Delete(ctx context.Context, user *appuser.User) error {
	result := repo.db.WithContext(ctx).Delete(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
