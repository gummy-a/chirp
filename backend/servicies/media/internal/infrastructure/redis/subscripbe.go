package redis

import (
	"encoding/json"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
)

func (h *QueueHandler) SubscribeStatus(jobId value_object.JobId) (<-chan entity.JobStatus, error) {
	statusCh := make(chan entity.JobStatus, 1)
	value, err := h.rdb.Get(h.ctx, jobId.String()).Result()
	if err == nil {
		var mediaId value_object.MediaId
		err = mediaId.ParseString(value)
		if err != nil {
			return nil, err
		}

		h.logger.Info("already encode finished.", slog.String("mediaId", mediaId.String()))

		statusCh <- entity.JobStatus{
			Status:  value_object.STATUS_COMPLETED,
			MediaId: &mediaId,
		}
		close(statusCh)
		return statusCh, nil
	}

	go func() {
		defer close(statusCh)

		sub := h.rdb.Subscribe(h.ctx, jobId.String())
		defer sub.Close()
		ch := sub.Channel()

		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}

				var status entity.JobStatus
				err = json.Unmarshal([]byte(msg.Payload), &status)
				if err != nil {
					h.logger.Error("failed to unmarshal status", slog.String("error", err.Error()))
					return
				}

				statusCh <- status
				if status.Status == value_object.STATUS_COMPLETED || status.Status == value_object.STATUS_ERROR {
					return
				}

			case <-h.ctx.Done():
				return
			}
		}
	}()

	return statusCh, nil
}
