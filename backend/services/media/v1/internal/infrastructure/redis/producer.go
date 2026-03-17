package redis

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"encoding/json"
	"log/slog"
)

func (h *QueueHandler) EnqueueJob(input *entity.EncodeJob) error {
	data, err := json.Marshal(input)
	if err != nil {
		h.logger.Error("json.Marshal failed", slog.String("error", err.Error()))
		return err
	}

	err = h.rdb.RPush(h.ctx, QueueName, data).Err()
	if err != nil {
		h.logger.Error("RPush failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}
