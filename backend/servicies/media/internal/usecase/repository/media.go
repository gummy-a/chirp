package repository

import (
	"context"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type MediaRepository interface {
	Save(ctx context.Context, files *[]entity.UploadedFileInfo, owner_account_id *domain.OwnerAccountId) error
	// Delete(ctx context.Context, media *entity.Media) error // TODO: imple Delete()
}
