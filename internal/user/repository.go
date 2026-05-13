package user

import (
	"context"
)

type Repository interface {
	List(ctx context.Context) []*User
	GetByID(ctx context.Context, id int) (*User, bool)
	GetByLogin(ctx context.Context, login string) (*User, bool)
	GetByTgId(ctx context.Context, id int64) (*User, bool)

	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
}
