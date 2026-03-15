package usecase

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/repository"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/usecase/port/queue"
)

type MediaUploadInput struct {
	Files          []value_object.UploadedFile
	OwnerAccountId value_object.OwnerAccountId
}

type MediaUploadOutput struct {
	OriginalFileName value_object.OriginalFileName
	JobId            value_object.JobId
}

type MediaUploadUseCase struct {
	queue queue.QueueHandler
	repo  repository.MediaRepository
}

func NewMediaUploadUseCase(r repository.MediaRepository, q queue.QueueHandler) *MediaUploadUseCase {
	return &MediaUploadUseCase{
		queue: q,
		repo:  r,
	}
}

func (u *MediaUploadUseCase) EnqueueEncode(input *MediaUploadInput) (*[]MediaUploadOutput, error) {
	var out []MediaUploadOutput
	for _, v := range input.Files {
		jobId := value_object.NewUniqueJobId()

		err := u.queue.EnqueueJob(&entity.EncodeJob{
			JobId: jobId,
			MediaInfo: value_object.MediaInfo{
				UploadedFile:   v,
				OwnerAccountId: input.OwnerAccountId,
			},
		})
		if err != nil {
			return nil, err
		}

		out = append(out, MediaUploadOutput{
			OriginalFileName: value_object.OriginalFileName(v.OriginalFileName),
			JobId:            jobId,
		})
	}

	return &out, nil
}
