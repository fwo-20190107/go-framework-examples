package router

import (
	"examples/pkg/adapter/framework/echo/handler"
	"examples/pkg/infra/framework/echo/middleware"
	"examples/pkg/infra/framework/echo/web"

	_ "examples/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetRoute(c *handler.AppContainer) *echo.Echo {
	r := echo.New()
	r.HTTPErrorHandler = web.ErrorHandler
	r.Use(middleware.Recover(), middleware.Logger.WithLogger())

	authorized := r.Group("/")
	authorized.Use(middleware.Auth.WithCheckToken)

	r.POST("/signup", web.EchoHandler(c.User.Signup).Exec)
	r.POST("/signin", web.EchoHandler(c.Session.Signin).Exec)
	authorized.DELETE("/signout", web.EchoHandler(c.Session.Signout).Exec)

	authorized.GET("/user", web.EchoHandler(c.User.GetAll).Exec)
	authorized.GET("/user/:user_id", web.EchoHandler(c.User.GetByID).Exec)
	authorized.PATCH("/user", web.EchoHandler(c.User.ModifyName).Exec)
	authorized.PATCH("/user/:user_id", web.EchoHandler(c.User.ModifyAuthority).Exec)

	// swagger UI
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	return r
}
