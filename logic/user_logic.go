package logic

import (
	"context"
	"database/sql"
	"errors"
	"examples/logic/repository"
	"examples/model"
)

type UserLogic interface {
	Login(ctx context.Context, loginID, password string) (*model.User, error)
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type userLogic struct {
	userRepository repository.UserRepository
}

func NewUserLogic(userRepository repository.UserRepository) *userLogic {
	return &userLogic{
		userRepository: userRepository,
	}
}

func (l *userLogic) Login(ctx context.Context, loginID, password string) (*model.User, error) {
	login, err := l.userRepository.GetLoginByID(ctx, loginID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("loginID is not exist")
		}
		return nil, err
	}

	if login.Password != password {
		return nil, errors.New("password does not match")
	}

	user, err := l.userRepository.GetUserByID(ctx, login.UserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (l *userLogic) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	user, err := l.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (l *userLogic) GetAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := l.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

var _ UserLogic = (*userLogic)(nil)
