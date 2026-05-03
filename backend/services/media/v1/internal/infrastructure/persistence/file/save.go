package file

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func sanitizeFileName(name string) string {
	name = filepath.Base(name)
	re := regexp.MustCompile(`[^a-zA-Z0-9\.\-_]`)
	return re.ReplaceAllString(name, "_")
}

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

func SaveFilesToTmp(files []*multipart.FileHeader) (*[]value_object.UploadedFile, error) {
	var ret []value_object.UploadedFile
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		safeName := sanitizeFileName(fileHeader.Filename)
		dst, err := os.CreateTemp("", safeName)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(dst, file)
		dst.Close() // defer dst.Close() actually closes when the loop reaches to the end

		if err != nil {
			return nil, err
		}

		mime, err := getMimeType(dst.Name())
		if err != nil {
			return nil, err
		}

		ret = append(ret, value_object.UploadedFile{
			OriginalFileName: value_object.OriginalFileName(fileHeader.Filename),
			RealPath:         value_object.RealPath(dst.Name()),
			FileUrl:          CreateUniqueFileUrl(),
			MimeType:         value_object.MimeType(mime),
		})
	}
	return &ret, nil
}
