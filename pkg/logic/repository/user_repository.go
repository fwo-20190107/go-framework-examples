package repository

import (
	"context"
	"examples/pkg/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, userID int) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Store(ctx context.Context, user *entity.User) (int64, error)
	ModifyAuthority(ctx context.Context, userID int, authority int8) error
	ModifyName(ctx context.Context, userID int, name string) error
}
