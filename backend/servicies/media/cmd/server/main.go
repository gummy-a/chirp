package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	controller "github.com/gummy_a/chirp/media/internal/adapter/controller"
	"github.com/gummy_a/chirp/media/internal/adapter/controller/router"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	repository "github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/infrastructure/redis"
	"github.com/gummy_a/chirp/media/internal/usecase"
	"github.com/joho/godotenv"
)

func setDefaultEnvironmentVariables() {
	env := os.Getenv("MEDIA_SERVICE_APP_ENV")
	if env == "development" {
		os.Setenv("MEDIA_SERVICE_PORT", "8081")
		os.Setenv("MEDIA_SERVICE_JWT_SECRET_KEY", "PSsDWRYMnGnLZpq1uq4Dd24WnGncTBkbtciiXzFNqGPHyJ") // must be same as auth service jwt-secret-key
		os.Setenv("MEDIA_SERVICE_ALLOW_ORIGIN", "http://localhost:3000")                            // DO NOT SET WILDCARD
		os.Setenv("MEDIA_SERVICE_DATABASE_URL", "postgres://postgres:password@localhost:5432/media_service?sslmode=disable")
		os.Setenv("MEDIA_SERVICE_REDIS_URL", "localhost:6379")
	} else {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf(".env not loaded. %v\n", err)
		}
	}
}

func checkEnvironmentVariables() {
	env := os.Getenv("MEDIA_SERVICE_APP_ENV")
	if env == "" {
		log.Fatal("MEDIA_SERVICE_APP_ENV environment variable is not set")
	}

	port := os.Getenv("MEDIA_SERVICE_PORT")
	if env == "production" && port == "" {
		log.Fatal("MEDIA_SERVICE_PORT environment variable is required in production")
	}

	jwtSecretKey := os.Getenv("MEDIA_SERVICE_JWT_SECRET_KEY")
	if env == "production" && jwtSecretKey == "" {
		log.Fatal("MEDIA_SERVICE_JWT_SECRET_KEY is not set")
	}

	allowOrigin := os.Getenv("MEDIA_SERVICE_ALLOW_ORIGIN")
	if env == "production" && allowOrigin == "" {
		log.Fatal("MEDIA_SERVICE_ALLOW_ORIGIN environment variable is required in production")
	}

	url := os.Getenv("MEDIA_SERVICE_DATABASE_URL")
	if env == "production" && url == "" {
		log.Fatal("MEDIA_SERVICE_DATABASE_URL is not set")
	}

	redis := os.Getenv("MEDIA_SERVICE_REDIS_URL")
	if env == "production" && redis == "" {
		log.Fatal("MEDIA_SERVICE_REDIS_URL is not set")
	}
}

func main() {
	setDefaultEnvironmentVariables()
	checkEnvironmentVariables()

	ctx := context.Background()

	// setup logger
	opts := &slog.HandlerOptions{AddSource: true}
	jsoncontroller := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsoncontroller)

	// Infrastructure layer: create DB pool
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Infrastructure layer: create database object
	queries := sqlc.New(pool)
	queue := redis.NewQueueHandler(ctx, *logger)

	// Repository layer: create repositories
	mediaRepository := repository.NewMediaRepository(pool, queries, logger)

	// UseCase layer: create use cases
	mediaControlUseCase := usecase.NewMediaUploadUseCase(mediaRepository, queue)

	// Adapter layer: create HTTP controllers and router
	mediaHandler := controller.NewUploadHandler(mediaControlUseCase, logger)
	router := router.NewAppRouter(mediaHandler, logger)

	//  Start HTTP server
	port := os.Getenv("MEDIA_SERVICE_PORT")
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
