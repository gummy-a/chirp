package fake

import (
	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"

	"github.com/google/uuid"
)

type MediaRepository struct {
	SaveMetaDataToDBCalls []struct {
		Metadata value_object.MetaData
		Job      entity.EncodeJob
	}
	SaveFileToStorageCalls []value_object.UploadedFile
}

func (f *MediaRepository) Save(media *entity.Media) (*value_object.MediaId, error) {
	// This is a fake, so we don't need to do anything here.
	id := value_object.MediaId(uuid.New())
	return &id, nil
}

func (f *MediaRepository) SaveFileToStorage(fname value_object.RealPath, url value_object.FileUrl) error {
	f.SaveFileToStorageCalls = append(f.SaveFileToStorageCalls, *file)
	return nil
}
