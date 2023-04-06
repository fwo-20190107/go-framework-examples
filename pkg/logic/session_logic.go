package logic

import (
	"context"
	"database/sql"
	"encoding/binary"
	"examples/pkg/code"
	"examples/pkg/errors"
	"examples/pkg/logger"
	"examples/pkg/logic/iodata"
	"examples/pkg/logic/repository"
	"examples/pkg/util"
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
	transaction       repository.Transaction
}

func NewSessionLogic(
	sessionRepository repository.SessionRepository,
	loginRepository repository.LoginRepository,
	transaction repository.Transaction,
) SessionLogic {
	return &sessionLogic{
		sonyflake:         *sonyflake.NewSonyflake(sonyflake.Settings{}),
		sessionRepository: sessionRepository,
		loginRepository:   loginRepository,
		transaction:       transaction,
	}
}

func (l *sessionLogic) Signin(ctx context.Context, input *iodata.SigninInput) (int, error) {
	login, err := l.loginRepository.GetByID(ctx, input.LoginID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.Errorf(code.CodeNotFound, err.Error())
		}
		return 0, err
	}

	// テストデータ投入コスト削減のため平文で登録/比較しています
	// 当然ダメなのでパスワードは最低限ハッシュ化して保存し、
	// 比較には bcrypt.CompareHashAndPassword() を使用すること
	if login.Password != input.Password {
		return 0, errors.Errorf(code.CodeBadRequest, "wrong loginID or password")
	}
	if _, err := l.transaction.Do(ctx, func(ctx context.Context) (interface{}, error) {
		if err := l.loginRepository.ModifyLastSigned(ctx, login.LoginID); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return 0, err
	}

	return login.UserID, nil
}

func (l *sessionLogic) Signout(ctx context.Context) {
	token, err := util.GetAccessToken(ctx)
	if err != nil {
		logger.L.Warn(err.Error())
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
		return "", errors.Wrap(code.CodeInternal, err)
	}

	bytesID := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(bytesID, id)

	token, err := bcrypt.GenerateFromPassword(bytesID, 4)
	if err != nil {
		return "", errors.Wrap(code.CodeInternal, err)
	}
	return string(token), nil
}

var _ SessionLogic = (*sessionLogic)(nil)
