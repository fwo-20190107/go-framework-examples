package session

import (
	"context"
	"encoding/binary"
	"examples/internal/http/entity/infra"
	"examples/internal/http/infra/cache"
	"examples/internal/http/interface/repository"
	"time"

	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

type sessionManager struct {
	sonyflake sonyflake.Sonyflake
	session   repository.LocalStore
}

func InitSessionManager() *sessionManager {
	return &sessionManager{
		sonyflake: *sonyflake.NewSonyflake(sonyflake.Settings{}),
		session:   cache.NewLocalStore(),
	}
}

func (m *sessionManager) NewToken() (string, error) {
	id, err := m.sonyflake.NextID()
	if err != nil {
		return "", err
	}

	bytesID := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(bytesID, id)

	token, err := bcrypt.GenerateFromPassword(bytesID, 4)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func (m *sessionManager) Load(ctx context.Context, token string) (int, bool) {
	v, ok := m.session.Get(ctx, token)
	if !ok {
		return 0, false
	}
	userID, ok := v.(int)
	if !ok {
		return 0, false
	}
	return userID, true
}

func (m *sessionManager) AddSession(ctx context.Context, token string, userID int) {
	m.session.Set(ctx, token, userID, 1*time.Hour)
}

func (m *sessionManager) RevokeSession(ctx context.Context, token string) {
	m.session.Drop(ctx, token)
}

var _ infra.SessionManage = (*sessionManager)(nil)
