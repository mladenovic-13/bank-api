package router

import "github.com/go-chi/chi/v5"

func (ctx *RouterContext) NewAccountRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", ctx.WithAuth(ctx.HandleGetAccounts))
	router.Post("/", ctx.WithAuth(ctx.HandleCreateAccount))

	router.Get("/{id}", ctx.WithAuth(ctx.HandleGetAccount))
	router.Delete("/{id}", ctx.WithAuth(ctx.HandleDeleteAccount))

	router.Post("/{number}/send", ctx.WithAuth(ctx.HandleSend))

	return router
}
