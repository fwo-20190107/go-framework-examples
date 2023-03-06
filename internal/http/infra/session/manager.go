package session

import (
	"encoding/binary"
	"examples/internal/http/entity/infra"
	"sync"

	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
)

type sessionManager struct {
	sonyflake sonyflake.Sonyflake
	session   sync.Map
}

func InitSessionManager() *sessionManager {
	return &sessionManager{
		sonyflake: *sonyflake.NewSonyflake(sonyflake.Settings{}),
		session:   sync.Map{},
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

func (m *sessionManager) Load(token string) (int, bool) {
	v, ok := m.session.Load(token)
	if !ok {
		return 0, false
	}
	userID, ok := v.(int)
	if !ok {
		return 0, false
	}
	return userID, true
}

func (m *sessionManager) AddSession(token string, userID int) {
	m.session.Store(token, userID)
}

func (m *sessionManager) RevokeSession(token string) {
	m.session.Delete(token)
}

var _ infra.SessionManage = (*sessionManager)(nil)
