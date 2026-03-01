package usecase

import (
	"context"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/usecase/repository"
)

type MediaUploadInput struct {
	Files          []entity.UploadedFileInfo
	OwnerAccountId domain.OwnerAccountId
}

type MediaUploadOutput struct {
	FileUrl  domain.FileUrl
	MimeType domain.MimeType
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

func (u *MediaControlUseCase) toMediaUploadOutput(in []entity.UploadedFileInfo) []MediaUploadOutput {
	var out []MediaUploadOutput
	for _, file := range in {
		out = append(out, MediaUploadOutput{
			FileUrl:  domain.FileUrl(file.FileUrl),
			MimeType: domain.MimeType(file.MimeType),
		})
	}
	return out
}

func (u *MediaControlUseCase) EnqueueEncode(ctx context.Context, input *MediaUploadInput) (*[]MediaUploadOutput, error) {
	for _, v := range input.Files {
		err := u.queue.EnqueueJob(&entity.EncodeJob{
			FileInfo: entity.UploadedFileInfo{
				UploadedFilePath: domain.UploadedFilePath(v.UploadedFilePath),
				FileUrl:          domain.FileUrl(v.FileUrl),
				MimeType:         domain.MimeType(v.MimeType),
			},
			OwnerAccountId: input.OwnerAccountId,
		})
		if err != nil {
			return nil, err
		}
	}

	out := u.toMediaUploadOutput(input.Files)
	return &out, nil
}
