package converter

import (
	"os"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/file"
)

func ToUploadedFileInfo(files []*os.File) ([]entity.UploadedFileInfo, error) {
	var entityFiles []entity.UploadedFileInfo
	for _, v := range files {
		mime, err := file.GetMimeType(v.Name())
		if err != nil {
			return nil, err
		}

		entityFiles = append(entityFiles, entity.UploadedFileInfo{
			OriginalFileName: value_object.OriginalFileName(v.Name()),
			FileUrl:          value_object.CreateAssetUrl(),
			MimeType:         value_object.MimeType(mime),
		})
	}
	return entityFiles, nil
}
