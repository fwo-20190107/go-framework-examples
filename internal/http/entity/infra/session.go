package infra

import "context"

type SessionManage interface {
	NewToken() (token string, err error)
	Load(token string) (userID int, ok bool)
	AddSession(token string, userID int)
	RevokeSession(ctx context.Context) (err error)
}
