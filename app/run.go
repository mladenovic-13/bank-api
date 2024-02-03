package app

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/mladenovic-13/bank-api/database"
)

func Run() error {
	err := godotenv.Load(".env")
	if err != nil {
		return errors.New("failed to load environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return errors.New("missing PORT environment variable")
	}

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return errors.New("missing DB_URL environment variable")
	}

	db, err := database.ConnectDB(connStr)

	if err != nil {
		return err
	}

	app, err := NewFiberApp(db)

	if err != nil {
		return err
	}

	app.Listen("localhost" + ":" + port)

	return nil
}
