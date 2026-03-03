package redis

import (
	"context"
	"log/slog"
	"os"

	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/redis/go-redis/v9"
)

type QueueHandler struct {
	rdb    redis.Client
	ctx    context.Context
	logger slog.Logger
	sql    sqlc.Queries
}

func NewQueueHandler(ctx context.Context, logger slog.Logger, sql sqlc.Queries) *QueueHandler {
	url := os.Getenv("MEDIA_SERVICE_REDIS_URL")
	return &QueueHandler{
		rdb: *redis.NewClient(&redis.Options{
			Addr: url,
		}),
		ctx:    ctx,
		logger: logger,
		sql: sql,
	}
}
