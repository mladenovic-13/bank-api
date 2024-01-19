package router

import "github.com/go-chi/chi/v5"

func NewAuthRouter(ctx *RouterContext) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/signin", ctx.HandleSignin)
	router.Post("/login", ctx.HandleLogin)
	router.Post("/logout", ctx.WithAuth(ctx.HandleLogout))

	return router
}
