package redis

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	MaxStreamLength    = 100
	ExpireStreamMinute = 5
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
		streamKey := value_object.NewStreamKey(&job.MediaInfo.OwnerAccountId)

		_, err = h.rdb.XAdd(h.ctx, &redis.XAddArgs{
			Stream: streamKey,
			MaxLen: MaxStreamLength,
			Approx: true,
			Values: map[string]interface{}{
				"msg":    fmt.Sprintf("start encoding %s ...", string(job.MediaInfo.UploadedFile.OriginalFileName)),
				"job_id": job.JobId.String(),
			},
		}).Result()
		if err != nil {
			h.logger.Error("XAdd failed", slog.String("error", err.Error()))
			continue
		}

		// set stream expire
		_, err = h.rdb.Expire(h.ctx, streamKey, ExpireStreamMinute*time.Minute).Result()
		if err != nil {
			h.logger.Error("Expire failed", slog.String("error", err.Error()))
			continue
		}

		mediaId, err := worker(&job)
		if err != nil {
			h.logger.Error("worker failed", slog.String("error", err.Error()))
			continue
		}

		_, err = h.rdb.XAdd(h.ctx, &redis.XAddArgs{
			Stream: streamKey,
			MaxLen: MaxStreamLength,
			Approx: true,
			Values: map[string]interface{}{
				"msg":      fmt.Sprintf("encode %s finished.", string(job.MediaInfo.UploadedFile.OriginalFileName)),
				"job_id":   job.JobId.String(),
				"media_id": mediaId.String(),
			},
		}).Result()
		if err != nil {
			h.logger.Error("XAdd failed", slog.String("error", err.Error()))
			continue
		}
	}
}
