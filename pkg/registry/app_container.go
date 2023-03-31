//go:build wireinject
// +build wireinject

package registry

import (
	"examples/pkg/adapter/handler"
	"examples/pkg/adapter/infra"
	"examples/pkg/adapter/repository"
	"examples/pkg/logic"

	"github.com/google/wire"
)

type AppContainer struct {
	User    handler.UserHandler
	Session handler.SessionHandler
}

func NewAppContainer(userHandler handler.UserHandler, sessionHandler handler.SessionHandler) *AppContainer {
	return &AppContainer{
		User:    userHandler,
		Session: sessionHandler,
	}
}

// wire用
func InitializeAppController(sqlh infra.SqlHandler, store infra.LocalStore) *AppContainer {
	wire.Build(
		NewAppContainer,
		handler.NewUserHandler,
		handler.NewSessionHandler,
		logic.NewUserLogic,
		logic.NewSessionLogic,
		repository.NewLoginRepository,
		repository.NewUserRepository,
		repository.NewSessionRepository,
	)
	return &AppContainer{}
}
