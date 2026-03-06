package redis

import (
	"encoding/json"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
)

const (
	QueueName = "encode_queue"
)

func (h *QueueHandler) EnqueueJob(input entity.EncodeJob) error {
	json, err := json.Marshal(input)
	if err != nil {
		h.logger.Error("json.Marshal failed", slog.String("error", err.Error()))
		return err
	}

	err = h.rdb.RPush(h.ctx, QueueName, json).Err()
	if err != nil {
		h.logger.Error("RPush failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}
