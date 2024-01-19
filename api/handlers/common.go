package handlers

import (
	"database/sql"

	"github.com/mladenovic-13/bank-api/internal/database"
)

type HandlerContext struct {
	DB      *sql.DB
	Queries *database.Queries
}

func NewHandlerContext(db *sql.DB, queries *database.Queries) *HandlerContext {
	return &HandlerContext{
		DB:      db,
		Queries: queries,
	}
}
