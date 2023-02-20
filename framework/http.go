package framework

import (
	"examples/infra/middleware"
	"examples/interface/handler"
	"examples/interface/repository"
	"examples/logic"
	"net/http"
)

func New(sqlh repository.SqlHandler) {
	u := handler.NewUserHandler(logic.NewUserLogic(repository.NewUserRepository(sqlh)))
	{
		http.Handle("/user", middleware.WithLogger(middleware.CheckToken(u)))
	}
}
