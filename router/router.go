package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mladenovic-13/bank-api/handlers"
)

func SetupRoutes(router *chi.Mux, ctx *handlers.RouterCtx) {
	router.Get("/healthz", ctx.HandleHealthz)

	v1Router := chi.NewRouter()
	v1Router.Post("/signin", ctx.HandleSignin)
	v1Router.Post("/login", ctx.HandleLogin)

	router.Mount("/v1", v1Router)
}
