package user

import (
	"context"
)

type Repository interface {
	List(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByTgId(ctx context.Context, id int64) (*User, error)

	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
}
