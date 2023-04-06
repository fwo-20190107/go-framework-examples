package router

import (
	"examples/pkg/adapter/framework/http/handler"
	"examples/pkg/infra/framework/http/middleware"
	"examples/pkg/infra/framework/http/web"
	"net/http"

	_ "examples/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRoute(c *handler.AppContainer) {
	defaultHandler := middleware.NewMiddlewareStack(middleware.LoggerMw.WithLogger, middleware.WithRecover)
	authHandler := defaultHandler.Append(middleware.AuthMw.WithCheckToken)

	// handler
	http.Handle("/signup", defaultHandler.Then(web.HttpHandler(c.User.Signup)))
	http.Handle("/signin", defaultHandler.Then(web.HttpHandler(c.Session.Signin)))
	http.Handle("/signout", authHandler.Then(web.HttpHandler(c.Session.Signout)))

	http.Handle("/user/", authHandler.Then(web.HttpHandler(c.User.HandleRoot)))

	// swagger UI
	http.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
}
