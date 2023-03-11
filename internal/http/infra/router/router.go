package router

import (
	"examples/internal/http/infra/middleware"
	"examples/internal/http/infra/web"
	"examples/internal/http/interface/handler"
	"examples/internal/http/interface/repository"
	"examples/internal/http/logic"
	"net/http"
)

func SetRoute(sqlh repository.SqlHandler) {
	u := handler.NewUserHandler(logic.NewUserLogic(repository.NewUserRepository(sqlh)))
	{
		http.Handle("/users", middleware.WithLogger(middleware.CheckToken(web.HttpHandler(u.HandleRoot))))
	}

	a := handler.NewSessionHandler(logic.NewUserLogic(repository.NewUserRepository(sqlh)))
	{
		http.Handle("/session", middleware.WithLogger(web.HttpHandler(a.HandleRoot)))
	}
}
