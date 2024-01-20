package router

import "github.com/go-chi/chi/v5"

func (ctx *RouterContext) NewAdminRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(ctx.WithAdmin)
	router.Post("/deposit/{number}", ctx.HandleDeposit)
	router.Post("/withdraw/{number}", ctx.HandleWithdraw)

	return router
}
