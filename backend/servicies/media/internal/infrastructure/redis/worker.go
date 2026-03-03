package redis

import (
	"encoding/json"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
)

type Worker func(entity.EncodeJob) error

func (h *QueueHandler) Worker(worker Worker) {
	for {
		result, err := h.rdb.BLPop(h.ctx, 0, QueueName).Result()
		if err != nil {
			h.logger.Error("BLPop failed", slog.String("error", err.Error()))
			continue
		}

		var job entity.EncodeJob
		err = json.Unmarshal([]byte(result[1]), &job) // key: result[0], value: result[1]
		if err != nil {
			h.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
			continue
		}

		err = worker(job)
		if err != nil {
			h.logger.Error("worker failed", slog.String("error", err.Error()))
			continue
		} else {
			h.logger.Info("worker success")
		}
	}
}
