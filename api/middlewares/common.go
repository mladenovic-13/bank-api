package middlewares

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/mladenovic-13/bank-api/internal/database"
)

type MiddlewareContext struct {
	DB      *sql.DB
	Queries *database.Queries
}

func NewMiddlewareContext(db *sql.DB, queries *database.Queries) *MiddlewareContext {
	return &MiddlewareContext{
		DB:      db,
		Queries: queries,
	}
}

func UseMiddlewares(
	router *chi.Mux,
	middlewares chi.Middlewares,
) {
	for _, middleware := range middlewares {
		router.Use(middleware)
	}
}
