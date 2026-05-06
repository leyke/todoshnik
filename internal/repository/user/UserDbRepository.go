package user

import (
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

func (repo *UserDbRepository) List() []*domain.User {
	var result []*domain.User
	repo.db.Find(result)

	return result
}

func (repo *UserDbRepository) GetByID(id int) (*domain.User, bool) {
	var item *domain.User
	result := repo.db.First(&item, id)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *UserDbRepository) GetByLogin(login string) (*domain.User, bool) {
	var item *domain.User
	result := repo.db.First(&item, "login = ?", login)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *UserDbRepository) GetUserByTgId(userTgId int64) (*domain.User, bool) {
	var item *domain.User
	result := repo.db.First(&item, "telegram_id = ?", userTgId)
	if result.Error != nil {
		return nil, false
	}

	return item, true
}

func (repo *UserDbRepository) Create(user *domain.User) (*domain.User, error) {
	result := repo.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *UserDbRepository) Update(user *domain.User) error {
	result := repo.db.Save(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *UserDbRepository) Delete(user *domain.User) error {
	result := repo.db.Delete(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
