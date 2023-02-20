package repository

import (
	"context"
	"examples/model"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
	GetLoginByID(ctx context.Context, loginID string) (*model.Login, error)
	ModifyAuthority(ctx context.Context, userID, authority int) error
}
