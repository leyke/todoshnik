package user

import "todoshnik/internal/domain"

type UserRepositoryInterface interface {
	List() []*domain.User
	GetByID(id int) (*domain.User, bool)
	GetByLogin(login string) (*domain.User, bool)
	GetUserByTgId(id int64) (*domain.User, bool)

	Create(user *domain.User) (*domain.User, error)
	Update(user *domain.User) error
	Delete(user *domain.User) error
}
