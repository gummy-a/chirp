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
	ExpireStreamMinute = 5 * time.Minute
)

type Worker func(*entity.EncodeJob) (*value_object.MediaId, error)

func (h *QueueHandler) startProcess(result []string) (*entity.EncodeJob, *string, error) {
	var job entity.EncodeJob
	err := json.Unmarshal([]byte(result[1]), &job) // key: result[0], value: result[1]
	if err != nil {
		h.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
		return nil, nil, err
	}
	streamKey := NewStreamKey(&job.MediaInfo.OwnerAccountId)
	pipe := h.rdb.Pipeline()

	_, err = pipe.XAdd(h.ctx, &redis.XAddArgs{
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
		return nil, nil, err
	}

	_, err = pipe.Expire(h.ctx, streamKey, ExpireStreamMinute).Result()
	if err != nil {
		h.logger.Error("Expire failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	_, err = pipe.Exec(h.ctx)
	if err != nil {
		h.logger.Error("Exec failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	return &job, &streamKey, nil
}

func (h *QueueHandler) endProcess(streamKey *string, job *entity.EncodeJob, mediaId *value_object.MediaId) error {
	_, err := h.rdb.XAdd(h.ctx, &redis.XAddArgs{
		Stream: *streamKey,
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
		return err
	}
	return nil
}

func (h *QueueHandler) Worker(worker func(*entity.EncodeJob) (*value_object.MediaId, error)) {
	for {
		result, err := h.rdb.BLPop(h.ctx, 0, QueueName).Result()
		if err != nil {
			h.logger.Error("BLPop failed", slog.String("error", err.Error()))
			continue
		}

		job, streamKey, err := h.startProcess(result)
		if err != nil {
			continue
		}

		mediaId, err := worker(job)
		if err != nil {
			h.logger.Error("worker failed", slog.String("error", err.Error()))
			continue
		}

		if err := h.endProcess(streamKey, job, mediaId); err != nil {
			continue
		}
	}
}
