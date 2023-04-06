package router

import (
	"examples/pkg/adapter/framework/gin/handler"
	"examples/pkg/infra/framework/gin/middleware"
	"examples/pkg/infra/framework/gin/web"

	_ "examples/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetRoute(c *handler.AppContainer) *gin.Engine {
	r := gin.New()
	r.Use(middleware.WithLogger(), gin.Recovery())

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMw.WithCheckToken())

	r.POST("/signup", web.GinHandler(c.User.Signup).Exec)
	r.POST("/signin", web.GinHandler(c.Session.Signin).Exec)
	authorized.DELETE("/signout", web.GinHandler(c.Session.Signout).Exec)

	authorized.GET("/user", web.GinHandler(c.User.GetAll).Exec)
	authorized.GET("/user/:user_id", web.GinHandler(c.User.GetByID).Exec)
	authorized.PATCH("/user", web.GinHandler(c.User.ModifyName).Exec)
	authorized.PATCH("/user/:user_id", web.GinHandler(c.User.ModifyAuthority).Exec)

	// swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)),
	)

	return r
}
