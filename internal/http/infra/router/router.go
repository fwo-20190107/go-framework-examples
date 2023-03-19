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
	// repository
	loginRepo := repository.NewLoginRepository(sqlh)
	userRepo := repository.NewUserRepository(sqlh)

	// logic
	loginLogic := logic.NewLoginLogic(userRepo, loginRepo)
	userLogic := logic.NewUserLogic(userRepo)

	users := handler.NewUserHandler(userLogic)
	{
		http.Handle("/users/", middleware.WithLogger(web.HttpHandler(middleware.CheckToken(users.HandleRoot))))
	}

	session := handler.NewSessionHandler(userLogic, loginLogic)
	{
		http.Handle("/session", middleware.WithLogger(web.HttpHandler(session.HandleRoot)))
	}
}
