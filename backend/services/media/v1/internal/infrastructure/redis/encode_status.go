package redis

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	blockSecond = 10
)

func (h *QueueHandler) Status(reqCtx context.Context, ownerAccountId *value_object.OwnerAccountId, workerFunc func(map[string]interface{})) error {
	// "$" marks events that are published only after the connection is established.
	//
	// NOTE:
	// Encoding may start before the client connects to SSE, so some events may be lost.
	lastId := "$"

	// to avoid infinite loop, set retry limit
	for range MaxStreamLength {
		select {
		case <-reqCtx.Done():
			return nil

		default:
			streams, err := h.rdb.XRead(h.ctx, &redis.XReadArgs{
				Streams: []string{NewStreamKey(ownerAccountId), lastId},
				Block:   blockSecond * time.Second,
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
				}
			}
		}
	}
	return nil
}
