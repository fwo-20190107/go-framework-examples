package infra

import "context"

type SessionManage interface {
	NewToken() (token string, err error)
	Load(ctx context.Context, token string) (userID int, ok bool)
	AddSession(ctx context.Context, token string, userID int)
	RevokeSession(ctx context.Context, token string)
}
