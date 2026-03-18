package redis

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	maxEncodeStreamLength = 1000
	maxSSEStreamLength    = 30
	expireStreamMinute    = 5 * time.Minute
	encodeStreamName      = "stream:encode:job"
	blockSecond           = 10 * time.Second
)

type Worker func(*entity.EncodeJob) (*value_object.MediaId, error)

func (h *QueueHandler) startProcess(result string) (*entity.EncodeJob, *string, error) {
	var job entity.EncodeJob
	if err := json.Unmarshal([]byte(result), &job); err != nil {
		h.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	sseStreamName := NewSSEStreamName(&job.MediaInfo.OwnerAccountId)
	pipe := h.rdb.Pipeline()

	_, err := pipe.XAdd(h.ctx, &redis.XAddArgs{
		Stream: sseStreamName,
		MaxLen: maxSSEStreamLength,
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

	_, err = pipe.Expire(h.ctx, sseStreamName, expireStreamMinute).Result()
	if err != nil {
		h.logger.Error("Expire failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	_, err = pipe.Exec(h.ctx)
	if err != nil {
		h.logger.Error("Exec failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	return &job, &sseStreamName, nil
}

func (h *QueueHandler) endProcess(sseStreamName *string, job *entity.EncodeJob, mediaId *value_object.MediaId) error {
	_, err := h.rdb.XAdd(h.ctx, &redis.XAddArgs{
		Stream: *sseStreamName,
		MaxLen: maxSSEStreamLength,
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

func (h *QueueHandler) execWorkerFunc(res []redis.XStream, worker Worker, groupName string) {
	for _, stream := range res {
		for _, msg := range stream.Messages {
			result, ok := msg.Values["payload"].(string)
			if !ok {
				h.logger.Error("type assertion failed")
				continue
			}

			job, sseStreamName, err := h.startProcess(result)
			if err != nil {
				continue
			}

			mediaId, err := worker(job)
			if err != nil {
				h.logger.Error("worker failed", slog.String("error", err.Error()))
				continue
			} else {
				if err := h.rdb.XAck(h.ctx, encodeStreamName, groupName, msg.ID).Err(); err != nil {
					h.logger.Error("XAck failed", slog.String("error", err.Error()))
					continue
				}
			}

			if err := h.endProcess(sseStreamName, job, mediaId); err != nil {
				continue
			}
		}
	}
}

func (h *QueueHandler) Worker(worker func(*entity.EncodeJob) (*value_object.MediaId, error)) {
	groupName := "group:worker"
	consumerName := "consumer:worker:" + uuid.NewString()
	err := h.rdb.XGroupCreateMkStream(h.ctx, encodeStreamName, groupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		h.logger.Error("XGroupCreateMkStream failed", slog.String("error", err.Error()))
		return
	}

	for {
		res, err := h.rdb.XReadGroup(h.ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{encodeStreamName, ">"},
			Block:    blockSecond,
			Count:    1,
		}).Result()

		if err == redis.Nil {
			continue
		}

		if err != nil {
			h.logger.Error("XReadGroup failed", slog.String("error", err.Error()))
			continue
		}

		h.execWorkerFunc(res, worker, groupName)
	}
}
