package sql

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"

	"github.com/mladenovic-13/bank-api/internal/database"
)

func NewPostgresStore() (*database.Queries, error) {
	url := os.Getenv("DB_URL")

	if url == "" {
		return nil, errors.New("failed to load db url env")
	}

	connection, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	db := database.New(connection)

	return db, nil
}
