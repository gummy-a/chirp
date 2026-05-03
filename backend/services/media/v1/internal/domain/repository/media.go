package repository

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
)

type MediaRepository interface {
	Save(media *entity.Media) (*value_object.MediaId, error)
	SaveFileToStorage(path value_object.RealPath, url value_object.FileUrl) error
	Delete(path value_object.RealPath) error
}
