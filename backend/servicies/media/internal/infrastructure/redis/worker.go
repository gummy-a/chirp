package redis

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type Worker func(entity.EncodeJob) (*value_object.MediaId, error)

func (h *QueueHandler) publishStatus(status entity.JobStatus, jobId string) error {
	msg, err := json.Marshal(status)
	if err != nil {
		return err
	}

	err = h.rdb.Publish(h.ctx, jobId, msg).Err()
	if err != nil {
		return err
	}

	return nil
}

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

		var status entity.JobStatus
		mediaId, err := worker(job)

		if err != nil {
			status.Status = value_object.STATUS_ERROR
			*status.Message = "failed encoding"
			h.logger.Error("worker failed", slog.String("error", err.Error()))
		} else {
			status.Status = value_object.STATUS_COMPLETED
			status.MediaId = mediaId
			*status.Message = "encode successed."

			// set completed before publish to avoid race condition
			err = h.rdb.Set(h.ctx, job.JobId.String(), mediaId.String(), 5*time.Minute).Err()
			if err != nil {
				h.logger.Error("Set failed", slog.String("error", err.Error()))
				continue
			}
		}

		err = h.publishStatus(status, job.JobId.String())
		if err != nil {
			h.logger.Error("publishStatus failed", slog.String("error", err.Error()))
			continue
		}

		h.logger.Info("Publish success", slog.String("job_id", job.JobId.String()))
	}
}
