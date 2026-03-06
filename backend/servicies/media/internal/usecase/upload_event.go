package usecase

import (
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/usecase/interface"
)

type UploadEventUseCase struct {
	queue interfaces.QueueHandler
	repo  repository.MediaRepository
}

func NewUploadEventUseCase(r repository.MediaRepository, q interfaces.QueueHandler) UploadEventUseCase {
	return UploadEventUseCase{
		queue: q,
		repo:  r,
	}
}

func (u *UploadEventUseCase) Execute(jobId value_object.JobId) (<-chan entity.JobStatus, error) {
	status, err := u.queue.SubscribeStatus(jobId)
	if err != nil {
		return nil, err
	}

	return status, nil
}
