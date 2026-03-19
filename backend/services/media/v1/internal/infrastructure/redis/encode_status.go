package redis

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	sseBlockSecond = 1 * time.Second
)

func (h *QueueHandler) deleteStream(streamName string) error {
	if err := h.rdb.Del(h.ctx, streamName).Err(); err != nil {
		h.logger.Error("Del failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (h *QueueHandler) Status(reqCtx context.Context, ownerAccountId *value_object.OwnerAccountId, workerFunc func(map[string]interface{})) error {
	// NOTE:
	// Encoding may start before the client connects to SSE, so some events may be lost.
	lastId := "$"

	streamName := NewSSEStreamName(ownerAccountId)
	defer h.deleteStream(streamName)

	for retryCount := 0; retryCount < 30; {
		select {
		case <-reqCtx.Done():
			return nil

		default:
			retryCount++
			streams, err := h.rdb.XRead(h.ctx, &redis.XReadArgs{
				Streams: []string{streamName, lastId},
				Block:   sseBlockSecond,
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
					workerFunc(msg.Values)
					lastId = msg.ID
					retryCount = 0
				}
			}
		}
	}
	return nil
}
