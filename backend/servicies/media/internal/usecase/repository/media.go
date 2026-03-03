package repository

import (
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

type MediaRepository interface {
	SaveMetaDataToDB(metadata entity.MetaData, job entity.EncodeJob) error
	SaveFileToStorage(url value_object.FileUrl) error
}
