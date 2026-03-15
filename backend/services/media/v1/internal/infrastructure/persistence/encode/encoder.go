package encode

import (
	"errors"
	"log/slog"

	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
)

type Encoder struct {
	logger *slog.Logger
}

func NewEncoder(logger *slog.Logger) *Encoder {
	return &Encoder{
		logger: logger,
	}
}

func (e *Encoder) Encode(job *entity.EncodeJob) (*value_object.MetaData, error) {
	var metadata *value_object.MetaData

	switch job.MediaInfo.UploadedFile.MimeType {
	// sometimes http.DetectContentType() can't detect mp4 and it returns mime type as application/octet-stream,
	// so we'll try to treat as octed-stream as mp4
	case "application/octet-stream":
		job.MediaInfo.UploadedFile.MimeType = "video/mp4"
		fallthrough
	case "video/mp4":
		m, err := e.encodeVideo(job)
		if err != nil {
			return nil, err
		}
		metadata = m

	case "image/png":
		fallthrough
	case "image/jpeg":
		fallthrough
	case "image/webp":
		m, err := e.encodeImage(job)
		if err != nil {
			return nil, err
		}
		metadata = m

	default:
		e.logger.Error("not allowed mime type", slog.String("mime type", string(job.MediaInfo.UploadedFile.MimeType)))
		return nil, errors.New("not allowed mime type")
	}

	return metadata, nil
}
