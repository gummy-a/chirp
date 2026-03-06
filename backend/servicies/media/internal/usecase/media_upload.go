package usecase

import (
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
	"github.com/gummy_a/chirp/media/internal/usecase/interface"
)

type MediaUploadInput struct {
	Files          []entity.UploadedFileInfo
	OwnerAccountId value_object.OwnerAccountId
}

type MediaUploadOutput struct {
	FileUrl  value_object.FileUrl
	MimeType value_object.MimeType
}

type MediaControlUseCase struct {
	queue interfaces.QueueHandler
	repo  repository.MediaRepository
}

func NewMediaUploadUseCase(r repository.MediaRepository, q interfaces.QueueHandler) MediaControlUseCase {
	return MediaControlUseCase{
		queue: q,
		repo:  r,
	}
}

func (u *MediaControlUseCase) toMediaUploadOutput(in []entity.UploadedFileInfo) []MediaUploadOutput {
	var out []MediaUploadOutput
	for _, file := range in {
		out = append(out, MediaUploadOutput{
			FileUrl:  value_object.FileUrl(file.FileUrl),
			MimeType: value_object.MimeType(file.MimeType),
		})
	}
	return out
}

func (u *MediaControlUseCase) EnqueueEncode(input MediaUploadInput) (*[]MediaUploadOutput, error) {
	for _, v := range input.Files {
		jobId, err := value_object.CreateJobId()
		if err != nil {
			return nil, err
		}

		err = u.queue.EnqueueJob(entity.EncodeJob{
			UploadedFileInfo: v,
			OwnerAccountId:   input.OwnerAccountId,
			JobId:            *jobId,
		})
		if err != nil {
			return nil, err
		}
	}

	out := u.toMediaUploadOutput(input.Files)
	return &out, nil
}
