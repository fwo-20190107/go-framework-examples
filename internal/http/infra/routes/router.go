package routes

import (
	"examples/internal/http/infra/middleware"
	"examples/internal/http/interface/handler"
	"examples/internal/http/interface/repository"
	"examples/internal/http/logic"
	"net/http"
)

func SetRoute(sqlh repository.SqlHandler) {
	u := handler.NewUserHandler(logic.NewUserLogic(repository.NewUserRepository(sqlh)))
	{
		http.Handle("/users", middleware.WithLogger(middleware.CheckToken(handler.AppHandler(u.HandleRoot))))
	}

	a := handler.NewSessionHandler(logic.NewUserLogic(repository.NewUserRepository(sqlh)))
	{
		http.Handle("/session", middleware.WithLogger(handler.AppHandler(a.HandleRoot)))
	}
}
