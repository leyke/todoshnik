package user

import (
	"context"
	"todoshnik/internal/domain"
)

type UserRepositoryInterface interface {
	List(ctx context.Context) []*domain.User
	GetByID(ctx context.Context, id int) (*domain.User, bool)
	GetByLogin(ctx context.Context, login string) (*domain.User, bool)
	GetUserByTgId(ctx context.Context, id int64) (*domain.User, bool)

	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, user *domain.User) error
}
