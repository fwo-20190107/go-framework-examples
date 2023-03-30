package repository

import (
	"context"
	"examples/pkg/entity"
)

type LoginRepository interface {
	GetByID(ctx context.Context, loginID string) (*entity.Login, error)
	GetByUserID(ctx context.Context, userID int) (*entity.Login, error)
	Store(ctx context.Context, login *entity.Login) error
	ModifyLastSigned(ctx context.Context, loginID string) error
}
