package logic

import (
	"context"
	"database/sql"
	"errors"
	"examples/internal/http/entity"
	"examples/internal/http/logic/repository"
	"examples/internal/http/registry"
	"examples/internal/http/util"
)

type UserLogic interface {
	Login(ctx context.Context, loginID, password string) (*entity.User, error)
	Logout(ctx context.Context)
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

type userLogic struct {
	userRepository repository.UserRepository
}

func NewUserLogic(userRepository repository.UserRepository) *userLogic {
	return &userLogic{
		userRepository: userRepository,
	}
}

func (l *userLogic) Login(ctx context.Context, loginID, password string) (*entity.User, error) {
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

	token, err := registry.SessionManager.NewToken()
	if err != nil {
		return nil, err
	}
	registry.SessionManager.AddSession(ctx, token, user.UserID)

	return user, nil
}

func (l *userLogic) Logout(ctx context.Context) {
	token, err := util.GetAccessToken(ctx)
	if err != nil {
		registry.Logger.Warn(ctx, err.Error())
		return
	}
	registry.SessionManager.RevokeSession(ctx, token)
}

func (l *userLogic) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	user, err := l.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (l *userLogic) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	users, err := l.userRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

var _ UserLogic = (*userLogic)(nil)
