//go:build wireinject
// +build wireinject

// Package registry is required `make wire`
package registry

import (
	"examples/pkg/adapter/handler"
	"examples/pkg/adapter/infra"
	"examples/pkg/adapter/repository"
	"examples/pkg/logic"

	"github.com/google/wire"
)

// wire用
func InitializeAppContainer(sqlh infra.SQLHandler, txh infra.TxHandler, store infra.LocalStore) *handler.AppContainer {
	wire.Build(
		handler.NewAppContainer,
		handler.NewUserHandler,
		handler.NewSessionHandler,
		logic.NewUserLogic,
		logic.NewSessionLogic,
		repository.NewLoginRepository,
		repository.NewUserRepository,
		repository.NewSessionRepository,
		repository.NewTransaction,
	)
	return &handler.AppContainer{}
}
