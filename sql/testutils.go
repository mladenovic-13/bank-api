package sql

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/mladenovic-13/bank-api/internal/database"
	"github.com/pressly/goose"
)

type TestStoreReturn struct {
	DB       *sql.DB
	Queries  *database.Queries
	Teardown func() error
}

func NewTestPostgresStore(connectionString string) (*TestStoreReturn, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	err = GooseUp(db)

	if err != nil {
		return nil, err
	}

	queries := database.New(db)

	teardown := func() error {
		return GooseDown(db)
	}

	return &TestStoreReturn{
		DB:       db,
		Queries:  queries,
		Teardown: teardown,
	}, nil
}

func GooseUp(db *sql.DB) error {
	err := goose.Up(db, "../../sql/schema")
	if err != nil {
		return err
	}

	return nil
}

func GooseDown(db *sql.DB) error {
	// TODO: Fix dirpath
	err := goose.Reset(db, "../../sql/schema")
	if err != nil {
		return err
	}

	return nil
}
