package registry

import (
	"examples/internal/http/entity/infra"
	"examples/internal/http/infra/log"
	"examples/internal/http/infra/session"
	"os"
)

var Logger infra.Logger

var SessionManager infra.SessionManage

func init() {
	Logger = log.InitLogger(os.Stdout)
	SessionManager = session.InitSessionManager()
}
