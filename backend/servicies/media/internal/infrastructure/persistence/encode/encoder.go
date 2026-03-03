package encode

import (
	"errors"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
)

type Encoder struct {
	logger slog.Logger
}

func NewEncoder(logger slog.Logger) Encoder {
	return Encoder{
		logger: logger,
	}
}

func (e *Encoder) Encode(job entity.EncodeJob) (*entity.MetaData, error) {
	var metadata *entity.MetaData

	switch job.FileInfo.MimeType {
	case "video/mp4":
		// TODO: implement video encode
		// output = string(job.FileInfo.UploadedFilePath) + ".encoded.mp4"
		// cmd := exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.UploadedFilePath), "-c:v", "libx264", "-crf", "25", "-c:a", "aac", output, "-y")

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
		e.logger.Error("not allowed mime type", slog.String("mime type", string(job.FileInfo.MimeType)))
		return nil, errors.New("not allowed mime type")
	}

	return metadata, nil
}
