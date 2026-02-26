package usecase

import (
	"context"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/usecase/repository"
)

type MediaUploadInput struct {
	Files          []entity.OriginalFileInfo
	OwnerAccountId domain.OwnerAccountId
}

type MediaUploadOutput struct {
	UnprocessedFileUrl domain.UnprocessedFileUrl
	FileType           domain.FileType
}

type QueueHandler interface {
	EnqueueJob(input *entity.EncodeJob) error
}

type MediaControlUseCase struct {
	media repository.MediaRepository
	queue QueueHandler
}

func NewMediaUploadUseCase(r repository.MediaRepository, q QueueHandler) *MediaControlUseCase {
	return &MediaControlUseCase{
		media: r,
		queue: q,
	}
}

func (u *MediaControlUseCase) toMediaUploadOutput(in []entity.OriginalFileInfo) []MediaUploadOutput {
	var out []MediaUploadOutput
	for _, file := range in {
		out = append(out, MediaUploadOutput{
			UnprocessedFileUrl: domain.UnprocessedFileUrl(file.UnprocessedFileUrl),
			FileType:           domain.FileType(file.FileType),
		})
	}
	return out
}

func (u *MediaControlUseCase) EnqueueEncode(ctx context.Context, input *MediaUploadInput) (*[]MediaUploadOutput, error) {
	for _, v := range input.Files {
		err := u.queue.EnqueueJob(&entity.EncodeJob{
			InputFile: domain.InputFile(v.UnprocessedFileUrl),
			MimeType:  domain.MimeType(v.FileType),
		})
		if err != nil {
			return nil, err
		}
	}
	
	out := u.toMediaUploadOutput(input.Files)
	return &out, nil
}
