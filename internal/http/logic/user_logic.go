package logic

import (
	"context"
	"examples/internal/http/entity"
	"examples/internal/http/logic/repository"
)

type UserLogic interface {
	GetByID(ctx context.Context, userID int) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
}

type userLogic struct {
	userRepository repository.UserRepository
}

func NewUserLogic(userRepository repository.UserRepository) *userLogic {
	return &userLogic{
		userRepository: userRepository,
	}
}

func (l *userLogic) GetByID(ctx context.Context, userID int) (*entity.User, error) {
	user, err := l.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (l *userLogic) GetAll(ctx context.Context) ([]entity.User, error) {
	users, err := l.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

var _ UserLogic = (*userLogic)(nil)
