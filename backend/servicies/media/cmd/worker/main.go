package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gummy_a/chirp/media/cmd"
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

	// exec encode server
	handler := redis.NewQueueHandler(ctx, *logger)
	handler.ExecuteJob()
}
