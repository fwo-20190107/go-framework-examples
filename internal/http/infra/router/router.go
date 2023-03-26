package router

import (
	"examples/internal/http/infra/cache"
	"examples/internal/http/infra/middleware"
	"examples/internal/http/infra/web"
	"examples/internal/http/interface/handler"
	"examples/internal/http/interface/infra"
	"examples/internal/http/interface/repository"
	"examples/internal/http/logic"
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
	loggerMid := middleware.NewLoggerMiddleware(os.Stdout)
	authMid := middleware.NewAuthMiddleware(sessionRepo)

	// ミドルウェアのラップが重なって見にくいので、デフォルトで使用するハンドラーとして先に定義しておく
	defaultHandler := func(fn infra.HttpHandler) http.HandlerFunc {
		return loggerMid.WithLogger(middleware.WithRecover(web.HttpHandler(fn)))
	}

	// handler
	users := handler.NewUserHandler(userLogic)
	{
		http.Handle("/signup", defaultHandler(users.Signup))
		http.Handle("/user/", defaultHandler(authMid.CheckToken(users.HandleRoot)))
	}

	session := handler.NewSessionHandler(userLogic, sessionLogic)
	{
		// サインアップ、サインインはトークン取得前なのでチェックを行わない
		http.Handle("/signin", defaultHandler(session.Signin))
		http.Handle("/signout", defaultHandler(authMid.CheckToken(session.Signout)))
	}

	// swagger UI
	http.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
}
