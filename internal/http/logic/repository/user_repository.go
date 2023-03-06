package repository

import (
	"context"
	"examples/internal/http/entity"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	GetLoginByID(ctx context.Context, loginID string) (*entity.Login, error)
	ModifyAuthority(ctx context.Context, userID, authority int) error
}
