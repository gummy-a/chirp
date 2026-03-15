package queue

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"context"
)

type QueueHandler interface {
	EnqueueJob(input *entity.EncodeJob) error
	Worker(func(*entity.EncodeJob) (*value_object.MediaId, error))
	Status(reqCtx context.Context, jobId *value_object.JobId, response func(map[string]interface{})) error
}
