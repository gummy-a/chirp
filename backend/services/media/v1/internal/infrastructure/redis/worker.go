package redis

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Worker func(*entity.EncodeJob) (*value_object.MediaId, error)

func (h *QueueHandler) Worker(worker func(*entity.EncodeJob) (*value_object.MediaId, error)) {
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

		_, err = h.rdb.XAdd(h.ctx, &redis.XAddArgs{
			Stream: job.JobId.String(),
			Values: map[string]interface{}{
				"msg":                "start encoding",
				"original_file_name": string(job.MediaInfo.UploadedFile.OriginalFileName),
			},
		}).Result()
		if err != nil {
			h.logger.Error("XAdd failed", slog.String("error", err.Error()))
			continue
		}

		mediaId, err := worker(&job)

		if err != nil {
			h.logger.Error("worker failed", slog.String("error", err.Error()))
			continue
		}

		_, err = h.rdb.XAdd(h.ctx, &redis.XAddArgs{
			Stream: job.JobId.String(),
			Values: map[string]interface{}{
				"msg":                "encode finished",
				"media_id":           mediaId.String(),
				"original_file_name": string(job.MediaInfo.UploadedFile.OriginalFileName),
			},
		}).Result()
		if err != nil {
			h.logger.Error("XAdd failed", slog.String("error", err.Error()))
			continue
		}

		h.logger.Info("Encode job success", slog.String("job_id", job.JobId.String()))
	}
}
