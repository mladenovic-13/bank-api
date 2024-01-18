package api

import "github.com/mladenovic-13/bank-api/internal/database"

type ServerContext struct {
	DB *database.Queries
}

func NewServerContext(db *database.Queries) *ServerContext {
	return &ServerContext{
		DB: db,
	}
}
