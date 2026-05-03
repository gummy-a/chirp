package main

import (
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/db"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/db/sqlc"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/encode"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/repository"
	"chirp/backend/services/media/v1/internal/infrastructure/redis"
	"context"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not loaded.%v\n", err)
	}

	ctx := context.Background()

	// setup logger
	opts := &slog.HandlerOptions{AddSource: true}
	jsoncontroller := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsoncontroller)

	// Infrastructure layer: create DB dbPool
	dbPool, err := db.NewConnectionPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create db pool: %v", err)
	}
	defer dbPool.Close()

	// Infrastructure layer: create database object
	sql := sqlc.New(dbPool)

	// Infrastructure layer: create queue
	queue := redis.NewQueueHandler(ctx, logger)

	// Infrastructure layer: create encoder
	encoder := encode.NewEncoder(logger)

	// Repository layer: create repositories
	s3, err := repository.NewS3Client(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	mediaRepository := repository.NewMediaRepository(logger, sql, s3, ctx)
	workerFunc := repository.NewSaveStrategy(encoder, *mediaRepository)

	log.Printf("Starting encode service...")
	for i := 0; i < runtime.NumCPU(); i++ {
		go queue.Worker(workerFunc)
	}

	select {}
}
