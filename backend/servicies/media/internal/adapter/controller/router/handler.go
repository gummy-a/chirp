package router

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/media/internal/adapter/controller"
	"github.com/gummy_a/chirp/media/internal/adapter/controller/helper"
	api "github.com/gummy_a/chirp/media/internal/infrastructure/http/openapi/media/go"
)

type AppHandler struct {
	uploadHandler *controller.UploadHandler
	logger        *slog.Logger
}

func NewAppRouter(upload *controller.UploadHandler, logger *slog.Logger) http.Handler {
	server := &AppHandler{
		uploadHandler: upload,
		logger:        logger,
	}
	DefaultAPIController := api.NewDefaultAPIController(server)
	router := api.NewRouter(DefaultAPIController)

	return helper.NewChain(router)
}

func (h *AppHandler) ApiMediaV1UploadPost(ctx context.Context, files []*os.File) (api.ImplResponse, error) {
	return h.uploadHandler.Upload(ctx, files)
}
