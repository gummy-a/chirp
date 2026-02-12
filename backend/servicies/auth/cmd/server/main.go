package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/auth/internal/adapter/handler"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/db/sqlc"
	repository "github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/repository/impl"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/register_account"
)

func setDefaultEnvironmentVariables() {
	// os.Setenv("AUTH_SERVICE_APP_ENV", "development")
	env := os.Getenv("AUTH_SERVICE_APP_ENV")
	if env == "development" {
		os.Setenv("AUTH_SERVICE_PORT", "8080")
		os.Setenv("AUTH_SERVICE_JWT_SECRET_KEY", "PSsDWRYMnGnLZpq1uq4Dd24WnGncTBkbtciiXzFNqGPHyJ")
		os.Setenv("AUTH_SERVICE_ALLOW_ORIGIN", "http://localhost:3000") // DO NOT SET WILDCARD
		os.Setenv("AUTH_SERVICE_DATABASE_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	}
}

func checkEnvironmentVariables() {
	env := os.Getenv("AUTH_SERVICE_APP_ENV")
	if env == "" {
		log.Fatal("AUTH_SERVICE_APP_ENV environment variable is not set")
	}

	port := os.Getenv("AUTH_SERVICE_PORT")
	if env == "production" && port == "" {
		log.Fatal("AUTH_SERVICE_PORT environment variable is required in production")
	}

	jwtSecretKey := os.Getenv("AUTH_SERVICE_JWT_SECRET_KEY")
	if env == "production" && jwtSecretKey == "" {
		log.Fatal("AUTH_SERVICE_JWT_SECRET_KEY is not set")
	}

	allowOrigin := os.Getenv("AUTH_SERVICE_ALLOW_ORIGIN")
	if env == "production" && allowOrigin == "" {
		log.Fatal("AUTH_SERVICE_ALLOW_ORIGIN environment variable is required in production")
	}

	url := os.Getenv("AUTH_SERVICE_DATABASE_URL")
	if env == "production" && url == "" {
		log.Fatal("AUTH_SERVICE_DATABASE_URL is not set")
	}
}

func main() {
	setDefaultEnvironmentVariables()
	checkEnvironmentVariables()

	ctx := context.Background()

	// setup logger
	opts := &slog.HandlerOptions{AddSource: true}
	jsonhandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsonhandler)

	// Infrastructure layer: create DB pool
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Infrastructure layer: create SQLC queries
	queries := sqlc.New(pool)

	// Repository layer: create repositories
	accountRepo := repository.NewAccountRepository(pool, queries, logger)
	temporaryAccountRepo := repository.NewTemporaryAccountRepository(queries, logger)
	registrationSenderRepo := repository.NewRegistrationSenderRepository(logger)

	// UseCase layer: create use cases
	SignupAccountUseCase := usecase.NewSignupAccountUseCase(accountRepo, temporaryAccountRepo)
	SignupTemporaryAccountUseCase := usecase.NewSignupTemporaryAccountUseCase(temporaryAccountRepo, registrationSenderRepo)

	// Adapter layer: create HTTP handlers and router
	router := handler.NewSignupRouter(SignupTemporaryAccountUseCase, SignupAccountUseCase, logger)

	//  Start HTTP server
	port := os.Getenv("AUTH_SERVICE_PORT")
	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
