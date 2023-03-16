package logic

import (
	"context"
	"database/sql"
	"errors"
	"examples/internal/http/logic/repository"
	"examples/internal/http/registry"
	"examples/internal/http/util"
	"examples/message"
)

type LoginLogic interface {
	Signin(ctx context.Context, loginID, password string) (int, error)
	PublishToken(ctx context.Context, userID int) (string, error)
	Signout(ctx context.Context)
}

type loginLogic struct {
	loginRepository repository.LoginRepository
}

func NewLoginLogic(userRepository repository.UserRepository, loginRepository repository.LoginRepository) *loginLogic {
	return &loginLogic{
		loginRepository: loginRepository,
	}
}

func (l *loginLogic) Signin(ctx context.Context, loginID, password string) (int, error) {
	login, err := l.loginRepository.GetByID(ctx, loginID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, message.ErrUserNotFound
		}
		return 0, err
	}

	if login.Password != password {
		return 0, message.ErrUnmatchPassword
	}

	if err := l.loginRepository.ModifyLastSigned(ctx, login.LoginID); err != nil {
		return 0, err
	}
	return login.UserID, nil
}

func (l *loginLogic) Signout(ctx context.Context) {
	token, err := util.GetAccessToken(ctx)
	if err != nil {
		registry.Logger.Warn(ctx, err.Error())
		return
	}
	registry.SessionManager.RevokeSession(ctx, token)
}

func (l *loginLogic) PublishToken(ctx context.Context, userID int) (string, error) {
	token, err := registry.SessionManager.NewToken()
	if err != nil {
		return "", err
	}
	registry.SessionManager.AddSession(ctx, token, userID)

	return token, nil
}

var _ LoginLogic = (*loginLogic)(nil)
