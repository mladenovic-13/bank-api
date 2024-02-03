package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mladenovic-13/bank-api/router"
	"gorm.io/gorm"
)

func NewFiberApp(db *gorm.DB) (*fiber.App, error) {
	app := fiber.New()

	app.Use(cors.New())

	router.SetupRoutes(
		app,
		router.NewRouterContext(db),
	)

	return app, nil
}
