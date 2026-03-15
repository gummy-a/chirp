package encode

import (
	"fmt"
	"log/slog"
	"os/exec"

	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/file"
)

const (
	VideoSizeLimit = "720"
)

func (e *Encoder) encodeVideo(job *entity.EncodeJob) (*value_object.MetaData, error) {
	var cmd *exec.Cmd
	encoded := string(job.MediaInfo.UploadedFile.RealPath) + ".encoded.mp4"
	thumbnail := string(job.MediaInfo.UploadedFile.RealPath) + ".thumbnail.webp"

	// encode
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.MediaInfo.UploadedFile.RealPath), "-vf", fmt.Sprintf("scale='if(gt(iw,ih),min(%s,iw),-2)':'if(gt(iw,ih),-2,min(%s,ih))'", VideoSizeLimit, VideoSizeLimit), "-c:v", "libx264", "-crf", "25", "-c:a", "aac", encoded, "-y")
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
