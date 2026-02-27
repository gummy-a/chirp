package controller

import (
	"context"
	"fmt"
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

	output, err := s.usecase.EnqueueEncode(ctx, &usecase.MediaUploadInput{
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

	var response []api.ApiMediaV1UploadPost200ResponseInner
	for _, v := range *output {
		response = append(response, api.ApiMediaV1UploadPost200ResponseInner{
			Message: fmt.Sprintf("upload %s accepted", v.UnprocessedFileUrl),
		})
	}

	return api.ImplResponse{Code: 200, Body: response}, nil
}
