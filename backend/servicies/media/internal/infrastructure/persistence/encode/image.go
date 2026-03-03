package encode

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

const (
	ImageSizeLimit = "2000"
	ThumbnailSize  = "512"
)

func (e *Encoder) getImageSize(path string) (*int, *int, error) {
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

	height, err := strconv.Atoi(strings.TrimSpace(splited[0]))
	if err != nil {
		e.logger.Error("Atoi failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	return &width, &height, nil
}

func (e *Encoder) encodeImage(job entity.EncodeJob) (*entity.MetaData, error) {
	var cmd *exec.Cmd
	encoded := string(job.FileInfo.OriginalFileName) + ".encoded.webp"
	thumbnail := string(job.FileInfo.OriginalFileName) + ".thumbnail.webp"

	// encode
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.OriginalFileName), "-vf", fmt.Sprintf("scale='if(gt(iw,ih),min(%s,iw),-1)':'if(gt(iw,ih),-1,min(%s,ih))'", ImageSizeLimit, ImageSizeLimit), "-q:v", "75", encoded, "-y")
	if err := cmd.Run(); err != nil {
		e.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
	}

	// create thumbnail
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.OriginalFileName), "-vf", fmt.Sprintf("scale='if(gt(a,1),-1,%s)':'if(gt(a,1),%s,-1)':flags=lanczos,crop=%s:%s", ThumbnailSize, ThumbnailSize, ThumbnailSize, ThumbnailSize), thumbnail, "-y")
	if err := cmd.Run(); err != nil {
		e.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
	}

	// get image size
	encodedWidth, encodedHeight, err := e.getImageSize(encoded)
	if err != nil {
		return nil, err
	}
	thumbnailWidth, thumbnailHeight, err := e.getImageSize(thumbnail)
	if err != nil {
		return nil, err
	}

	ret := entity.MetaData{
		Thumbnail: entity.Asset{
			URL:         string(value_object.CreateAssetUrl()),
			Width:       *thumbnailWidth,
			Height:      *thumbnailHeight,
			Type:        "image/webp",
			VideoLength: nil,
		},
		Encoded: []entity.Asset{
			{
				URL:         string(value_object.CreateAssetUrl()),
				Width:       *encodedWidth,
				Height:      *encodedHeight,
				Type:        string(job.FileInfo.MimeType),
				VideoLength: nil,
			},
		},
	}

	return &ret, nil
}
