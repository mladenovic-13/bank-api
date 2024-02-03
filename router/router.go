package router

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, routerCtx *RouterContext) {
	rc := NewRouteContext(routerCtx)

	// BASE PATH
	v1api := app.Group("/api/v1")

	// MIDDLEWARES
	v1api.Use(logger.New())
	v1api.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
	}))

	// AUTH
	auth := v1api.Group("/auth")

	auth.Post("/signin", rc.HandleSignin)
	auth.Post("/login", rc.HandleLogin)
	auth.Post("/logout", rc.WithAuth, rc.HandleLogout)

	// ADMIN
	admin := v1api.Group("/admin", rc.WithAuth, rc.WithAdmin)

	admin.Get("/request", rc.HandleGetRequests)
	admin.Post("/request/:id/process", rc.HandleProcessRequest)
	admin.Post("/create", rc.HandleCreateAccount)
	admin.Post("/delete/:id", rc.HandleDeleteAccount) // TODO: implement
	admin.Post("/deposit/:id", rc.HandleDeposit)
	admin.Post("/withdraw/:id", rc.HandleWithdraw)

	// ACCOUNT
	account := v1api.Group("/account", rc.WithAuth)

	account.Get("/", rc.HandleGetAccounts)
	account.Post("/", rc.HandleRequestAccount)
	account.Post("/", rc.HandleCreateAccount)
	account.Get("/:id", rc.HandleGetAccount)
	account.Post("/:id/deposit", rc.HandleRequestDeposit)
	account.Post("/:id/withdraw", rc.HandleRequestWithdraw)
}
