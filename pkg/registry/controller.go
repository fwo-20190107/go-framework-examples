package registry

import "examples/pkg/adapter/handler"

type DIContainer struct {
	User    handler.UserHandler
	Session handler.SessionHandler
}
