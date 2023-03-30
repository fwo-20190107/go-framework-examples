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

func InitializeSessionHandler(sqlh infra.SqlHandler) *handler.SessionHandler {
	wire.Build(
		handler.NewSessionHandler,
		logic.NewUserLogic,
		repository.NewLoginRepository,
		repository.NewUserRepository,
	)
	return &handler.SessionHandler{}
}
