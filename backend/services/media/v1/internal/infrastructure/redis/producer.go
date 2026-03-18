package redis

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func (h *QueueHandler) EnqueueJob(input *entity.EncodeJob) error {
	data, err := json.Marshal(input)
	if err != nil {
		h.logger.Error("json.Marshal failed", slog.String("error", err.Error()))
		return err
	}

	err = h.rdb.XAdd(h.ctx, &redis.XAddArgs{
		Stream: encodeStreamName,
		MaxLen: maxEncodeStreamLength,
		Approx: true,
		Values: map[string]interface{}{
			"payload": data,
		},
	}).Err()
	if err != nil {
		h.logger.Error("XAdd failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}
