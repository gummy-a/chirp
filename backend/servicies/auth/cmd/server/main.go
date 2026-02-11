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

func main() {
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
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
