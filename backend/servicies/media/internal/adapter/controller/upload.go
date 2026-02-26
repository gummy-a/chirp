package controller

import (
	"context"
	"log/slog"
	"os"

	"github.com/gummy_a/chirp/media/internal/adapter/converter"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/http/middleware"
	api "github.com/gummy_a/chirp/media/internal/infrastructure/http/openapi/media/go"
	"github.com/gummy_a/chirp/media/internal/usecase"
)

type UploadHandler struct {
	usecase *usecase.MediaControlUseCase
	logger  *slog.Logger
}

func NewUploadHandler(usecase *usecase.MediaControlUseCase, logger *slog.Logger) *UploadHandler {
	return &UploadHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (s *UploadHandler) Upload(ctx context.Context, files []*os.File) (api.ImplResponse, error) {
	ownerAccountId, ok := ctx.Value(middleware.OwnerAccountIdKey).(domain.OwnerAccountId)
	if !ok {
		s.logger.Error("Failed to get accountId from context")
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to get accountId",
		}}, nil
	}

	originamFileInfo, err := converter.ToOriginalFileInfo(files)
	if err != nil {
		s.logger.Error("Failed ToOriginalFileInfo()", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to EnqueueEncode",
		}}, nil
	}

	_, err = s.usecase.EnqueueEncode(ctx, &usecase.MediaUploadInput{
		Files:          *originamFileInfo,
		OwnerAccountId: ownerAccountId,
	})
	if err != nil {
		s.logger.Error("Failed to EnqueueEncode", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to EnqueueEncode",
		}}, nil
	}

	// TODO: response 200 OK
	return api.ImplResponse{Code: 200, Body: []api.ApiMediaV1UploadPost200ResponseInner{
		{Message: "success"},
	}}, nil
}
