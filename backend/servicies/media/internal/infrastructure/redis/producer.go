package redis

import (
	"encoding/json"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
)

const (
	QueueName = "encode_queue"
)

func (h *QueueHandler) EnqueueJob(input *entity.EncodeJob) error {
	job := entity.EncodeJob{
		InputFile: domain.InputFile(input.InputFile),
		MimeType:  domain.MimeType(input.MimeType),
	}

	json, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = h.rdb.RPush(h.ctx, QueueName, json).Err()
	if err != nil {
		return err
	}

	return nil
}
