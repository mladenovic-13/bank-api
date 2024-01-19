package api

import (
	"database/sql"

	"github.com/mladenovic-13/bank-api/internal/database"
)

type ServerContext struct {
	DB      *sql.DB
	Queries *database.Queries
}

func NewServerContext(
	db *sql.DB,
	queries *database.Queries,
) *ServerContext {
	return &ServerContext{
		DB:      db,
		Queries: queries,
	}
}
