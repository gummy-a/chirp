package mediaController

import (
	"chirp/backend/services/media/v1/internal/usecase"
	"log/slog"
)

type MediaControllerUseCases struct {
	Upload usecase.MediaUploadUseCase
	Status usecase.MediaEncodeStatusUseCase
}

type MediaController struct {
	UseCases *MediaControllerUseCases
	Logger   *slog.Logger
}
