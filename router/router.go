package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/api/handlers"
	"github.com/mladenovic-13/bank-api/api/middlewares"
)

func SetupRoutes(router *chi.Mux, serverContext *api.ServerContext) {
	ctx := newRouterContext(serverContext)

	router.Get("/healthz", ctx.HandleHealthz)

	v1Router := chi.NewRouter()

	v1Router.Post("/signin", ctx.HandleSignin)
	v1Router.Post("/login", ctx.HandleLogin)
	v1Router.Get("/account", ctx.WithAuth(ctx.HandleGetAccounts))

	router.Mount("/v1", v1Router)
}

type RouterContext struct {
	*handlers.HandlerContext
	*middlewares.MiddlewareContext
}

func newRouterContext(serverContext *api.ServerContext) *RouterContext {
	return &RouterContext{
		HandlerContext:    handlers.NewHandlerContext(serverContext.DB),
		MiddlewareContext: middlewares.NewMiddlewareContext(serverContext.DB),
	}
}
