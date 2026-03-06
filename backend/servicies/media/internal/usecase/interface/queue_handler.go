package interfaces

import (
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type QueueHandler interface {
	EnqueueJob(input entity.EncodeJob) error
	SubscribeStatus(jobId value_object.JobId) (<-chan entity.JobStatus, error)
}
