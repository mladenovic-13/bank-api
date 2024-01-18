package middlewares

import "github.com/mladenovic-13/bank-api/internal/database"

type MiddlewareContext struct {
	DB *database.Queries
}

func NewMiddlewareContext(db *database.Queries) *MiddlewareContext {
	return &MiddlewareContext{
		DB: db,
	}
}
