package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/media/cmd"
	controller "github.com/gummy_a/chirp/media/internal/adapter/controller"
	"github.com/gummy_a/chirp/media/internal/adapter/controller/router"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	repository "github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/infrastructure/redis"
	"github.com/gummy_a/chirp/media/internal/usecase"
)

func main() {
	cmd.SetDefaultEnvironmentVariables()
	cmd.CheckEnvironmentVariables()
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
