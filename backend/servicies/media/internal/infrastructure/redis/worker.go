package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	ImageSizeLimit = "2000"
	ThumbnailSize  = "512"
)

func (h *QueueHandler) getImageSize(path string) (*int, *int, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", path)
	outputBytes, err := cmd.Output()
	if err != nil {
		h.logger.Error("ffprobe failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	outputString := string(outputBytes)
	splited := strings.Split(outputString, ",")

	width, err := strconv.Atoi(strings.TrimSpace(splited[0]))
	if err != nil {
		h.logger.Error("Atoi failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	height, err := strconv.Atoi(strings.TrimSpace(splited[0]))
	if err != nil {
		h.logger.Error("Atoi failed", slog.String("error", err.Error()))
		return nil, nil, err
	}

	return &width, &height, nil
}

func (h *QueueHandler) encodeImage(job entity.EncodeJob) (*entity.MetaData, error) {
	var cmd *exec.Cmd
	output := string(job.FileInfo.OriginalFileName) + ".encoded.webp"
	thumbnail := string(job.FileInfo.OriginalFileName) + ".thumbnail.webp"

	// encode
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.OriginalFileName), "-vf", fmt.Sprintf("scale='if(gt(iw,ih),min(%s,iw),-1)':'if(gt(iw,ih),-1,min(%s,ih))'", ImageSizeLimit, ImageSizeLimit), "-q:v", "75", output, "-y")
	if err := cmd.Run(); err != nil {
		h.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
	}

	// create thumbnail
	cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.OriginalFileName), "-vf", fmt.Sprintf("scale='if(gt(a,1),-1,%s)':'if(gt(a,1),%s,-1)':flags=lanczos,crop=%s:%s", ThumbnailSize, ThumbnailSize, ThumbnailSize, ThumbnailSize), thumbnail, "-y")
	if err := cmd.Run(); err != nil {
		h.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
	}

	// get image size
	encodedWidth, encodedHeight, err := h.getImageSize(output)
	if err != nil {
		return nil, err
	}
	thumbnailWidth, thumbnailHeight, err := h.getImageSize(thumbnail)
	if err != nil {
		return nil, err
	}

	ret := entity.MetaData{
		Thumbnail: entity.Asset{
			URL:         string(domain.CreateAssetUrl()),
			Width:       *thumbnailWidth,
			Height:      *thumbnailHeight,
			Type:        "image/webp",
			VideoLength: nil,
		},
		Encoded: []entity.Asset{
			{
				URL:         string(domain.CreateAssetUrl()),
				Width:       *encodedWidth,
				Height:      *encodedHeight,
				Type:        string(job.FileInfo.MimeType),
				VideoLength: nil,
			},
		},
	}

	return &ret, nil
}

func (h *QueueHandler) encode(job entity.EncodeJob) (*entity.MetaData, error) {
	var metadata *entity.MetaData

	switch job.FileInfo.MimeType {
	case "video/mp4":
		// output = string(job.FileInfo.UploadedFilePath) + ".encoded.mp4"
		// cmd := exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.UploadedFilePath), "-c:v", "libx264", "-crf", "25", "-c:a", "aac", output, "-y")

	case "image/png":
		fallthrough
	case "image/jpeg":
		fallthrough
	case "image/webp":
		m, err := h.encodeImage(job)
		if err != nil {
			return nil, err
		}
		metadata = m

	default:
		h.logger.Error("not allowed mime type", slog.String("mime type", string(job.FileInfo.MimeType)))
		return nil, errors.New("not allowed mime type")
	}

	return metadata, nil
}

func (h *QueueHandler) insert(metadata entity.MetaData, job entity.EncodeJob) error {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(job.OwnerAccountId),
		Valid: true,
	}

	meta, err := json.Marshal(metadata)
	if err != nil {
		h.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
		return err
	}

	_, err = h.sql.InsertMedia(h.ctx, sqlc.InsertMediaParams{
		OwnerAccountID:     pgtypeUUID,
		MimeType:           string(job.FileInfo.MimeType),
		OriginalFileName:   string(job.FileInfo.OriginalFileName),
		UnprocessedFileUrl: string(job.FileInfo.FileUrl),
		Metadata:           meta,
	})
	if err != nil {
		h.logger.Error("InsertMedia failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (h *QueueHandler) worker() {
	for {
		result, err := h.rdb.BLPop(h.ctx, 0, QueueName).Result()
		if err != nil {
			h.logger.Error("BLPop failed", slog.String("error", err.Error()))
			continue
		}

		var job entity.EncodeJob
		err = json.Unmarshal([]byte(result[1]), &job) // key: result[0], value: result[1]
		if err != nil {
			h.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
			continue
		}

		metadata, err := h.encode(job)
		if err != nil {
			continue
		}

		err = h.insert(*metadata, job)
		if err != nil {
			continue
		}

		h.logger.Info("Encoding finished successfully: ", slog.String("input", string(job.FileInfo.OriginalFileName)))
	}
}

func (h *QueueHandler) ExecuteJob() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go h.worker()
	}

	select {}
}
