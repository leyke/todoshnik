package user

import "todoshnik/internal/domain"

type UserRepositoryInreface interface {
	List() []*domain.User
	GetByID(id int) (*domain.User, error)
	GetUserByTgId(id int64) (*domain.User, error)

	Create(user *domain.User) (*domain.User, error)
	Update(user *domain.User) error
	Delete(user *domain.User) error
}
