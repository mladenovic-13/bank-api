package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/mladenovic-13/bank-api/api"
	"github.com/mladenovic-13/bank-api/api/router"
	"github.com/mladenovic-13/bank-api/sql"
)

func SetupTestServe(url string) (*chi.Mux, func() error, error) {
	store, err := sql.NewTestPostgresStore(url)

	if err != nil {
		return nil, nil, err
	}

	apiContext := api.NewServerContext(store.DB, store.Queries)

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

	return app, store.Teardown, nil
}
