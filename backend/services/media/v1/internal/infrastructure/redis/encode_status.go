package redis

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func (h *QueueHandler) Status(reqCtx context.Context, jobId *value_object.JobId, response func(map[string]interface{})) error {
	for {
		select {
		case <-reqCtx.Done():
			return nil

		default:
			lastId := "$"
			streams, err := h.rdb.XRead(h.ctx, &redis.XReadArgs{
				Streams: []string{jobId.String(), lastId},
				Block:   30 * time.Second,
			}).Result()

			if err == redis.Nil {
				continue
			}

			if err != nil {
				h.logger.Error("XRead failed", slog.String("error", err.Error()))
				return err
			}

			for _, s := range streams {
				for _, msg := range s.Messages {
					response(msg.Values)
					lastId = msg.ID
				}
			}
		}
	}
}
