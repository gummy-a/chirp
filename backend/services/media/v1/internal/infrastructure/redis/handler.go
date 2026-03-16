package redis

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
)

const (
	QueueName = "encode_queue"
)

type QueueHandler struct {
	rdb    *redis.Client
	ctx    context.Context
	logger *slog.Logger
}

func NewQueueHandler(ctx context.Context, logger *slog.Logger) *QueueHandler {
	url := os.Getenv("MEDIA_SERVICE_REDIS_URL")
	return &QueueHandler{
		rdb: redis.NewClient(&redis.Options{
			Addr: url,
		}),
		ctx:    ctx,
		logger: logger,
	}
}

func NewStreamKey(ownerAccountId *value_object.OwnerAccountId) string {
	return "stream:encode:" + ownerAccountId.String()
}
