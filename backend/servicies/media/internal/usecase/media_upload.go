package usecase

import (
	"context"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
	repository "github.com/gummy_a/chirp/media/internal/infrastructure/persistence/repository/impl"
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
	queue QueueHandler
	repo  repository.MediaRepository
}

func NewMediaUploadUseCase(r repository.MediaRepository, q QueueHandler) *MediaControlUseCase {
	return &MediaControlUseCase{
		queue: q,
		repo:  r,
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

func (u *MediaControlUseCase) UploadToStorage(ctx context.Context, input *MediaUploadInput) (*[]entity.EncodeJob, error) {
	// TODO: implement this
	u.repo.Save(ctx, &input.Files, &input.OwnerAccountId)
	return nil, nil
}

func (u *MediaControlUseCase) EnqueueEncode(ctx context.Context, input *MediaUploadInput) (*[]MediaUploadOutput, error) {
	for _, v := range input.Files {
		err := u.queue.EnqueueJob(&entity.EncodeJob{
			FileInfo: entity.UploadedFileInfo{
				OriginalFileName: domain.OriginalFileName(v.OriginalFileName),
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
