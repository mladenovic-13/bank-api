package app

import (
	"errors"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Run() error {
	err := godotenv.Load(".env")

	if err != nil {
		panic("failed to load .env")
	}

	port := os.Getenv("PORT")

	if port == "" {
		return errors.New("failed to load port env variable")
	}

	url := os.Getenv("DB_URL")

	if url == "" {
		return errors.New("failed to load DB_URL env")
	}

	app, err := SetupServer(url)

	if err != nil {
		return err
	}

	err = http.ListenAndServe(":"+port, app)

	if err != nil {
		return err
	}

	return nil
}
