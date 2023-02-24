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
		http.Handle("/users", middleware.WithLogger(middleware.CheckToken(handler.AppHandler(u.HandleRoot))))
	}

	a := handler.NewSessionHandler(logic.NewUserLogic(repository.NewUserRepository(sqlh)))
	{
		http.Handle("/session", middleware.WithLogger(handler.AppHandler(a.HandleRoot)))
	}
}
