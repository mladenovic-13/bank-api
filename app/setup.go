package app

import (
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mladenovic-13/bank-api/db"
	"github.com/mladenovic-13/bank-api/handlers"
	"github.com/mladenovic-13/bank-api/router"
)

func SetupAndRunApp() error {
	err := godotenv.Load(".env")

	if err != nil {
		panic("failed to load .env")
	}

	port := os.Getenv("PORT")

	if port == "" {
		return errors.New("failed to load env variable")
	}

	db, err := db.NewPostgresStore()

	if err != nil {
		panic("failed to create connection to database")
	}

	routerCtx := handlers.NewRouterCtx(db)

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

	router.SetupRoutes(app, routerCtx)

	err = http.ListenAndServe(":"+port, app)

	if err != nil {
		return err
	}

	return nil
}
