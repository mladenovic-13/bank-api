package handlers

import (
	"github.com/mladenovic-13/bank-api/internal/database"
)

type HandlerContext struct {
	DB *database.Queries
}

func NewHandlerContext(db *database.Queries) *HandlerContext {
	return &HandlerContext{
		DB: db,
	}
}
