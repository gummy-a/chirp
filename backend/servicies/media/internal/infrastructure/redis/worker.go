package redis

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os/exec"
	"runtime"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
)

func (h *QueueHandler) encode(job entity.EncodeJob) (*exec.Cmd, *string, error) {
	var cmd *exec.Cmd
	var output string

	switch job.FileInfo.MimeType {
	case "video/mp4":
		output = string(job.FileInfo.UploadedFilePath) + ".encoded.mp4"
		cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.UploadedFilePath), "-c:v", "libx264", "-crf", "25", "-c:a", "aac", output, "-y")

	case "image/png":
		fallthrough
	case "image/jpeg":
		fallthrough
	case "image/webp":
		output = string(job.FileInfo.UploadedFilePath) + ".encoded.webp"
		cmd = exec.Command("ffmpeg", "-threads", "1", "-i", string(job.FileInfo.UploadedFilePath), "-q:v", "75", output, "-y")

	default:
		h.logger.Error("not allowed mime type", slog.String("mime type", string(job.FileInfo.MimeType)))
		return nil, nil, errors.New("not allowed mime type")
	}

	return cmd, &output, nil
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

		cmd, output, err := h.encode(job)
		if err != nil {
			continue
		}

		if err := cmd.Run(); err != nil {
			h.logger.Error("ffmpeg failed", slog.String("error", err.Error()))
			continue
		}

		h.logger.Info("Encoding finished successfully: ", slog.String("input", string(job.FileInfo.UploadedFilePath)), slog.String("output", *output))
	}
}

func (h *QueueHandler) ExecuteJob() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go h.worker()
	}

	select {}
}
