package infra

type SessionManage interface {
	NewToken() (token string, err error)
	Load(token string) (userID int, ok bool)
	AddSession(token string, userID int)
	RevokeSession(token string)
}
