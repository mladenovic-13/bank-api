package app

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/api/router"
	"github.com/mladenovic-13/bank-api/sql"
)

func SetupServer(url string) (*chi.Mux, error) {
	if url == "" {
		return nil, errors.New("failed to load db url env")
	}

	db, queries, err := sql.NewPostgresStore(url)

	if err != nil {
		return nil, err
	}

	apiContext := api.NewServerContext(db, queries)

	app := chi.NewRouter()

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}

	app.Use(cors.Handler(corsOptions))

	router.SetupRoutes(app, apiContext)

	return app, nil
}
