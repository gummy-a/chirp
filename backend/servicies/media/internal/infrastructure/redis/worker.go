package redis

import (
	"encoding/json"
	"log/slog"
	"os/exec"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
)

func (h *QueueHandler) ExecuteJob() {
	for {
		result, err := h.rdb.BLPop(h.ctx, 0, QueueName).Result()
		if err != nil {
			h.logger.Error("BLPop failed", slog.String("error", err.Error()))
			continue
		}

		var job entity.EncodeJob
		json.Unmarshal([]byte(result[1]), &job) // key: result[0], value: result[1]

		var cmd *exec.Cmd
		var output string

		switch job.FileInfo.MimeType {
		case "video/mp4":
			output = string(job.FileInfo.UploadedFilePath) + ".encoded.mp4"
			cmd = exec.Command("ffmpeg", "-i", string(job.FileInfo.UploadedFilePath), "-c:v", "libx264", "-crf", "25", "-c:a", "aac", output, "-y")

		case "image/png":
			fallthrough
		case "image/jpeg":
			fallthrough
		case "image/webp":
			output = string(job.FileInfo.UploadedFilePath) + ".encoded.webp"
			cmd = exec.Command("ffmpeg", "-i", string(job.FileInfo.UploadedFilePath), "-q:v", "75", output, "-y")

		default:
			h.logger.Error("not allowed mime type", slog.String("mime type", string(job.FileInfo.MimeType)))
			continue
		}

		if err := cmd.Run(); err != nil {
			h.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
		} else {
			h.logger.Info("Encoding finished successfully: ", slog.String("input", string(job.FileInfo.UploadedFilePath)), slog.String("output", output))
		}
	}
}
