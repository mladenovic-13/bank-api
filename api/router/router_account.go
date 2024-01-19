package router

import "github.com/go-chi/chi/v5"

func NewAccountRouter(ctx *RouterContext) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", ctx.WithAuth(ctx.HandleGetAccounts))
	router.Post("/", ctx.WithAuth(ctx.HandleCreateAccount))
	router.Get("/{id}", ctx.WithAuth(ctx.HandleGetAccount))
	router.Delete("/{id}", ctx.WithAuth(ctx.HandleDeleteAccount))
	router.Post("/{number}/deposit", ctx.WithAuth(ctx.HandleDeposit))
	router.Post("/{number}/withdraw", ctx.WithAuth(ctx.HandleWithdraw))

	return router
}
