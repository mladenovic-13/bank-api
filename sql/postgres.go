package sql

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"

	"github.com/mladenovic-13/bank-api/internal/database"
)

func NewPostgresStore() (*sql.DB, *database.Queries, error) {
	url := os.Getenv("DB_URL")

	if url == "" {
		return nil, nil, errors.New("failed to load db url env")
	}

	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, nil, err
	}

	queries := database.New(db)

	return db, queries, nil
}
