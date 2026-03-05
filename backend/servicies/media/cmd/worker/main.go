package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/gummy_a/chirp/media/cmd"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/encode"
	repository "github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/infrastructure/redis"
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

	// Infrastructure layer: create encoder
	encoder := encode.NewEncoder(*logger)

	// Infrastructure layer: create queue
	queue := redis.NewQueueHandler(ctx, *logger)

	// Repository layer: create repositories
	s3, err := repository.NewS3Client(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	mediaRepository := repository.NewMediaRepository(*logger, *sql, *s3, ctx)
	workerFunc := repository.NewSaveStrategy(encoder, mediaRepository)

	log.Printf("Starting encode service...")
	for i := 0; i < runtime.NumCPU(); i++ {
		go queue.Worker(workerFunc)
	}

	select {}
}
