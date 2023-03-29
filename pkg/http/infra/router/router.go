package router

import (
	"examples/pkg/http/adapter/handler"
	"examples/pkg/http/adapter/infra"
	"examples/pkg/http/adapter/repository"
	"examples/pkg/http/infra/cache"
	"examples/pkg/http/infra/middleware"
	"examples/pkg/http/infra/web"
	"examples/pkg/http/logic"
	"net/http"
	"os"

	_ "examples/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRoute(sqlh infra.SqlHandler) {
	// repository
	loginRepo := repository.NewLoginRepository(sqlh)
	userRepo := repository.NewUserRepository(sqlh)
	sessionRepo := repository.NewSessionRepository(cache.NewLocalStore())

	// logic
	userLogic := logic.NewUserLogic(userRepo, loginRepo)
	sessionLogic := logic.NewSessionLogic(sessionRepo, loginRepo)

	// middleware
	loggerMw := middleware.NewLoggerMiddleware(os.Stdout)
	authMw := middleware.NewAuthMiddleware(sessionRepo)

	defaultHandler := middleware.NewMiddlewareStack(loggerMw.WithLogger, middleware.WithRecover)
	authHandler := defaultHandler.Append(authMw.WithCheckToken)

	// handler
	users := handler.NewUserHandler(userLogic)
	{
		http.Handle("/signup", defaultHandler.Then(web.HttpHandler(users.Signup)))
		http.Handle("/user/", authHandler.Then(web.HttpHandler(users.HandleRoot)))
	}

	session := handler.NewSessionHandler(userLogic, sessionLogic)
	{
		// サインインはトークン取得前なのでチェックを行わない
		http.Handle("/signin", defaultHandler.Then(web.HttpHandler(session.Signin)))
		http.Handle("/signout", authHandler.Then(web.HttpHandler(session.Signout)))
	}

	// swagger UI
	http.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
}
