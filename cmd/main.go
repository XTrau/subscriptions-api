package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/database"
	"subscriptions-api/internal/handlers"
	"subscriptions-api/internal/middlewares"
	"subscriptions-api/internal/repositories"
	"subscriptions-api/internal/usecases"

	_ "subscriptions-api/docs"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscriptions API
// @version 1.0
// @host     localhost:8080
// @BasePath /

func main() {
	postgres, err := database.NewPostresDB(config.AppConfig)

	if err != nil {
		log.Fatal("Error on creating Postres connection pool.", err)
	} else {
		log.Println("Postgres connected!")
	}

	err = database.RunMigrations(config.AppConfig)

	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations.")
		} else {
			log.Fatal("Error on migrations running.", err)
		}
	} else {
		log.Println("Migrations complete!")
	}

	textHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(textHandler)

	repo := repositories.NewSubscriptionsPostgresRepository(postgres)
	ucases := usecases.NewSubscriptionUseCases(repo)
	sr := handlers.NewSubscriptionsRoutes(ucases, logger)

	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware(logger))

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	sr.RegisterRoutes(r)

	log.Println("Server started!")

	http.ListenAndServe(":8080", r)
}
