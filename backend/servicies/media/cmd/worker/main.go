package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gummy_a/chirp/media/cmd"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/gummy_a/chirp/media/internal/infrastructure/redis"
)

func main() {
	cmd.SetDefaultEnvironmentVariables()
	cmd.CheckEnvironmentVariables()
	ctx := context.Background()

	// Infrastructure layer: create DB pool
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Infrastructure layer: create database object
	queries := sqlc.New(pool)

	// setup logger
	opts := &slog.HandlerOptions{AddSource: true}
	jsoncontroller := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsoncontroller)

	// exec encode server
	handler := redis.NewQueueHandler(ctx, *logger, *queries)
	log.Printf("Starting encode service...")
	handler.ExecuteJob()
}
