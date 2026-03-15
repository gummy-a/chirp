package encode

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"

	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/file"
)

const (
	ImageSizeLimit = "2000"
	ThumbnailSize  = "512"
)

func (e *Encoder) getFileResolution(path string) (*int, *int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", path)
	outputBytes, err := cmd.Output()
	if err != nil {
		e.logger.Error("ffprobe failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	outputString := string(outputBytes)
	splited := strings.Split(outputString, ",")

	width, err := strconv.Atoi(strings.TrimSpace(splited[0]))
	if err != nil {
		e.logger.Error("Atoi failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	height, err := strconv.Atoi(strings.TrimSpace(splited[1]))
	if err != nil {
		e.logger.Error("Atoi failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	return &width, &height, nil
}

func (e *Encoder) encodeImage(job *entity.EncodeJob) (*value_object.MetaData, error) {
	var cmd *exec.Cmd
	encoded := string(job.MediaInfo.UploadedFile.RealPath) + ".encoded.webp"
	thumbnail := string(job.MediaInfo.UploadedFile.RealPath) + ".thumbnail.webp"

	// encode
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.MediaInfo.UploadedFile.RealPath), "-vf", fmt.Sprintf("scale='if(gt(iw,ih),min(%s,iw),-1)':'if(gt(iw,ih),-1,min(%s,ih))'", ImageSizeLimit, ImageSizeLimit), "-q:v", "75", encoded, "-y")
	if err := cmd.Run(); err != nil {
		e.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
	}

	// create thumbnail
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", encoded, "-vf", fmt.Sprintf("scale='if(gt(a,1),-1,%s)':'if(gt(a,1),%s,-1)':flags=lanczos,crop=%s:%s", ThumbnailSize, ThumbnailSize, ThumbnailSize, ThumbnailSize), thumbnail, "-y")
	if err := cmd.Run(); err != nil {
		e.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
	}

	// get image size
	encodedWidth, encodedHeight, err := e.getFileResolution(encoded)
	if err != nil {
		return nil, err
	}
	thumbnailWidth, thumbnailHeight, err := e.getFileResolution(thumbnail)
	if err != nil {
		return nil, err
	}

	ret := value_object.MetaData{
		Thumbnail: value_object.Asset{
			EncodedFile: value_object.EncodedFile{
				RealPath: value_object.RealPath(thumbnail),
				FileUrl:  value_object.FileUrl(file.CreateUniqueFileUrl()),
				MimeType: "image/webp",
			},
			Width:       value_object.Width(*thumbnailWidth),
			Height:      value_object.Height(*thumbnailHeight),
			VideoLength: nil,
		},
		Encoded: []value_object.Asset{
			{
				EncodedFile: value_object.EncodedFile{
					RealPath: value_object.RealPath(encoded),
					FileUrl:  value_object.FileUrl(file.CreateUniqueFileUrl()),
					MimeType: value_object.MimeType(job.MediaInfo.UploadedFile.MimeType),
				},
				Width:       value_object.Width(*encodedWidth),
				Height:      value_object.Height(*encodedHeight),
				VideoLength: nil,
			},
		},
	}

	return &ret, nil
}
