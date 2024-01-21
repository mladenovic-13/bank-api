package sql

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/mladenovic-13/bank-api/internal/database"
)

func NewPostgresStore(connectionString string) (*sql.DB, *database.Queries, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, nil, err
	}

	queries := database.New(db)

	return db, queries, nil
}
