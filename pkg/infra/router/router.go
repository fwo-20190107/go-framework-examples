package router

import (
	"examples/pkg/adapter/handler"
	"examples/pkg/infra/middleware"
	"examples/pkg/infra/web"
	"net/http"

	_ "examples/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRoute(c *handler.AppContainer) {
	defaultHandler := middleware.NewMiddlewareStack(middleware.LoggerMw.WithLogger, middleware.WithRecover)
	authHandler := defaultHandler.Append(middleware.AuthMw.WithCheckToken)

	// handler
	http.Handle("/signup", defaultHandler.Then(web.HttpHandler(c.User.Signup)))
	http.Handle("/user/", authHandler.Then(web.HttpHandler(c.User.HandleRoot)))
	http.Handle("/signin", defaultHandler.Then(web.HttpHandler(c.Session.Signin)))
	http.Handle("/signout", authHandler.Then(web.HttpHandler(c.Session.Signout)))

	// swagger UI
	http.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
}
