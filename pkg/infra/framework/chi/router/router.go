package router

import (
	"examples/pkg/adapter/framework/chi/handler"
	"examples/pkg/infra/framework/chi/middleware"
	"examples/pkg/infra/framework/chi/web"
	"os"

	_ "examples/docs"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetRoute(c *handler.AppContainer) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(middleware.WithLogger(os.Stdout), chimw.Recoverer)

		r.Group(func(r chi.Router) {
			r.Post("/signin", web.ChiHandler(c.Session.Signin).Exec)
			r.Post("/signup", web.ChiHandler(c.User.Signup).Exec)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth.WithCheckToken)

			r.Delete("/signout", web.ChiHandler(c.Session.Signout).Exec)

			r.Route("/user", func(r chi.Router) {
				r.Get("/", web.ChiHandler(c.User.GetAll).Exec)
				r.Get("/{user_id:^[1-9][0-9]*$}", web.ChiHandler(c.User.GetByID).Exec)
				r.Patch("/", web.ChiHandler(c.User.ModifyName).Exec)
				r.Patch("/{user_id:^[1-9][0-9]*$}", web.ChiHandler(c.User.ModifyAuthority).Exec)
			})
		})
	})

	// swagger UI
	r.Mount("/swagger", httpSwagger.WrapHandler)
	return r
}
