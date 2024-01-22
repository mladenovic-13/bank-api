package router

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/api/handlers"
	"github.com/mladenovic-13/bank-api/api/middlewares"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func SetupRoutes(router *chi.Mux, serverContext *api.ServerContext) {
	regularMiddlewares := chi.Middlewares{
		middlewares.RequestLogger,
		middlewares.RecoverPanic,
	}
	middlewares.UseMiddlewares(router, regularMiddlewares)

	ctx := newRouterContext(serverContext)

	router.Get("/healthz", ctx.HandleHealthz)

	router.Get("/docs/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	v1Router := chi.NewRouter()

	v1Router.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/docs/swagger.json"),
	))

	v1Router.Mount("/", ctx.NewAuthRouter())
	v1Router.Mount("/account", ctx.NewAccountRouter())
	v1Router.Mount("/admin", ctx.NewAdminRouter())

	router.Mount("/v1", v1Router)
}

type RouterContext struct {
	*handlers.HandlerContext
	*middlewares.MiddlewareContext
}

func newRouterContext(serverContext *api.ServerContext) *RouterContext {
	return &RouterContext{
		HandlerContext: handlers.NewHandlerContext(
			serverContext.DB,
			serverContext.Queries,
		),
		MiddlewareContext: middlewares.NewMiddlewareContext(
			serverContext.DB,
			serverContext.Queries,
		),
	}
}
