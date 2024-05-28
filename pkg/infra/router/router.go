package router

import (
	"examples/pkg/adapter/handler"
	"examples/pkg/infra/middleware"
	"examples/pkg/infra/web"

	_ "examples/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetRoute(c *handler.AppContainer) *gin.Engine {
	r := gin.New()
	r.Use(middleware.WithLogger(), gin.Recovery())

	authorized := r.Group("/")
	authorized.Use(middleware.Auth.WithCheckToken())

	r.POST("/signup", web.Handler(c.User.Signup).Exec)
	r.POST("/signin", web.Handler(c.Session.Signin).Exec)
	authorized.DELETE("/signout", web.Handler(c.Session.Signout).Exec)

	authorized.GET("/user", web.Handler(c.User.GetAll).Exec)
	authorized.GET("/user/:user_id", web.Handler(c.User.GetByID).Exec)
	authorized.PATCH("/user", web.Handler(c.User.ModifyName).Exec)
	authorized.PATCH("/user/:user_id", web.Handler(c.User.ModifyAuthority).Exec)

	// swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)),
	)

	return r
}
