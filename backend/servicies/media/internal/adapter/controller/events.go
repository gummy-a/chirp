package controller

import (
	"context"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	api "github.com/gummy_a/chirp/media/internal/infrastructure/http/openapi/media/go"
)

func (s *UploadHandler) Events(ctx context.Context, jobId string) (api.ImplResponse, error) {
	var job_id value_object.JobId
	err := job_id.ParseString(jobId)
	if err != nil {
		s.logger.Error("Failed to parse jobId", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to parse jobId",
		}}, nil
	}

	statusCh, err := s.m.Execute(job_id)
	if err != nil {
		s.logger.Error("Failed to execute usecase", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to parse jobId",
		}}, nil
	}

	eventCh := make(chan api.UploadProgressEvent)
	go func() {
		defer close(eventCh)
		for status := range statusCh {
			event := api.UploadProgressEvent{
				Status: string(status.Status),
			}

			if status.MediaId != nil {
				event.MediaId = status.MediaId.String()
			}
			if status.Message != nil {
				event.Message = string(*status.Message)
			}

			eventCh <- event
		}
	}()

	return api.ImplResponse{
		Code: 200,
		Body: api.SSEEventStream{Events: eventCh},
	}, nil
}
