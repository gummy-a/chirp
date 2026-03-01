package converter

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
)

func getMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func ToUploadedFileInfo(files []*os.File) (*[]entity.UploadedFileInfo, error) {
	var entityFiles []entity.UploadedFileInfo
	for _, v := range files {
		mime, err := getMimeType(v.Name())
		if err != nil {
			return nil, err
		}

		uuid := uuid.NewString()

		entityFiles = append(entityFiles, entity.UploadedFileInfo{
			UploadedFilePath: domain.UploadedFilePath(v.Name()),
			FileUrl:          domain.FileUrl("/assets/media/" + uuid),
			MimeType:         domain.MimeType(mime),
		})
	}
	return &entityFiles, nil
}
