package redis

import (
	"encoding/json"
	"log/slog"
	"os/exec"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	domain "github.com/gummy_a/chirp/media/internal/domain/value_object"
)

func (h *QueueHandler) ExecuteJob(inputFile domain.InputFile) {
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

		switch job.MimeType {
		case "video/mp4":
			output = string(job.InputFile) + ".encoded.mp4"
			cmd = exec.Command("ffmpeg", "-i", string(job.InputFile), "-c:v", "libx264", "-crf", "25", "-c:a", "aac", output, "-y")

		case "image/png":
			fallthrough
		case "image/jpeg":
			fallthrough
		case "image/webp":
			output = string(job.InputFile) + ".encoded.webp"
			cmd = exec.Command("ffmpeg", "-i", string(job.InputFile), "-q:v", "75", output, "-y")

		default:
			h.logger.Error("not allowed mime type", slog.String("mime type", string(job.MimeType)))
			continue
		}

		if err := cmd.Run(); err != nil {
			h.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
		} else {
			h.logger.Info("Encoding finished successfully: ", slog.String("input", string(job.InputFile)), slog.String("output", output))
		}
	}
}
