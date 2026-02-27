package redis

import (
	"encoding/json"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
)

const (
	QueueName = "encode_queue"
)

func (h *QueueHandler) EnqueueJob(input *entity.EncodeJob) error {
	json, err := json.Marshal(*input)
	if err != nil {
		return err
	}

	err = h.rdb.RPush(h.ctx, QueueName, json).Err()
	if err != nil {
		return err
	}

	return nil
}
