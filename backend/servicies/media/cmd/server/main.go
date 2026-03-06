package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/media/cmd"
	"github.com/gummy_a/chirp/media/internal/adapter/controller"
	"github.com/gummy_a/chirp/media/internal/adapter/controller/router"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
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

	// Infrastructure layer: create DB dbPool
	dbPool, err := db.NewConnectionPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbPool.Close()

	// Infrastructure layer: create database object
	sql := sqlc.New(dbPool)
	queue := redis.NewQueueHandler(ctx, *logger)

	// Repository layer: create repositories
	s3, err := repository.NewS3Client(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	mediaRepository := repository.NewMediaRepository(*logger, *sql, *s3, ctx)

	// UseCase layer: create use cases
	mediaControlUseCase := usecase.NewMediaUploadUseCase(mediaRepository, queue)
	uploadEventUseCase := usecase.NewUploadEventUseCase(mediaRepository, queue)


	// Adapter layer: create HTTP controllers and router
	mediaHandler := controller.NewUploadHandler(mediaControlUseCase, uploadEventUseCase, *logger)
	router := router.NewAppRouter(mediaHandler, *logger)

	//  Start HTTP server
	port := os.Getenv("MEDIA_SERVICE_PORT")
	log.Printf("Starting media service server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
