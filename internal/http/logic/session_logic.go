package logic

import (
	"context"
	"database/sql"
	"encoding/binary"
	"examples/code"
	"examples/errors"
	"examples/internal/http/logic/iodata"
	"examples/internal/http/logic/repository"
	"examples/internal/http/registry"
	"examples/internal/http/util"
	"time"

	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

type SessionLogic interface {
	Signin(ctx context.Context, input *iodata.SigninInput) (int, error)
	Signout(ctx context.Context)
	Start(ctx context.Context, userID int) (string, error)
}

type sessionLogic struct {
	sonyflake         sonyflake.Sonyflake
	sessionRepository repository.SessionRepository
	loginRepository   repository.LoginRepository
}

func NewSessionLogic(sessionRepository repository.SessionRepository, loginRepository repository.LoginRepository) *sessionLogic {
	return &sessionLogic{
		sonyflake:         *sonyflake.NewSonyflake(sonyflake.Settings{}),
		sessionRepository: sessionRepository,
		loginRepository:   loginRepository,
	}
}

func (l *sessionLogic) Signin(ctx context.Context, input *iodata.SigninInput) (int, error) {
	login, err := l.loginRepository.GetByID(ctx, input.LoginID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.Errorf(code.ErrNotFound, err.Error())
		}
		return 0, err
	}

	if login.Password != input.Password {
		return 0, errors.Errorf(code.ErrBadRequest, "wrong loginID or password")
	}

	if err := l.loginRepository.ModifyLastSigned(ctx, login.LoginID); err != nil {
		return 0, err
	}
	return login.UserID, nil
}

func (l *sessionLogic) Signout(ctx context.Context) {
	token, err := util.GetAccessToken(ctx)
	if err != nil {
		registry.Logger.Warn(ctx, err.Error())
		return
	}
	l.sessionRepository.Drop(ctx, token)
}

func (l *sessionLogic) Start(ctx context.Context, userID int) (string, error) {
	token, err := l.publishToken(ctx)
	if err != nil {
		return "", err
	}

	l.sessionRepository.Set(ctx, token, userID, 1*time.Hour)
	return token, nil
}

func (m *sessionLogic) publishToken(ctx context.Context) (string, error) {
	id, err := m.sonyflake.NextID()
	if err != nil {
		return "", errors.Wrap(code.ErrInternal, err)
	}

	bytesID := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(bytesID, id)

	token, err := bcrypt.GenerateFromPassword(bytesID, 4)
	if err != nil {
		return "", errors.Wrap(code.ErrInternal, err)
	}
	return string(token), nil
}

var _ SessionLogic = (*sessionLogic)(nil)
