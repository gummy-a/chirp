package handler

import (
	"chirp/backend/router"
	"chirp/backend/services/media/v1/internal/adapter/controller"
	"chirp/backend/services/media/v1/internal/infrastructure/http/middleware"
	"log/slog"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func applyMiddleware(handlers []router.Handler, mws ...Middleware) {
	for i := range handlers {
		var final http.Handler = handlers[i].Handler
		for j := len(mws) - 1; j >= 0; j-- {
			final = mws[j](final)
		}
		handlers[i].Handler = final.ServeHTTP
	}
}

func NewMediaHandler(uc *mediaController.MediaControllerUseCases, logger *slog.Logger) []router.Handler {
	controller := mediaController.MediaController{
		UseCases: uc,
		Logger:   logger,
	}
	mediaController := []router.Handler{
		{
			// Upload is a handler for file upload. This saves files to /tmp and enqueue encode jobs.
			// Input: files (binary) from multipart form values.
			// Success: 200 OK with JSON body {"original_file_name": "name", "job_id": "..."}
			// Failure: 400 Bad Request for invalid input, 500 Internal Server Error for failed enqueue.
			Path:    "POST /api/media/v1/upload/",
			Handler: controller.Upload,
		},
		{
			// Status is a handler for file encoding job.
			// this is SSE endpoint, connection is keep alived and responses encoding progress.
			// Input: jobId (string) from path.
			// Success: 200 OK with JSON body {"original_file_name": "name", "job_id": "..."}
			// Failure: 400 Bad Request for invalid input, 500 Internal Server Error for failed encoding.
			Path:    "GET /api/media/v1/upload/status/{jobId}/",
			Handler: controller.Status,
		},
	}

	applyMiddleware(mediaController, middleware.JwtMiddleware)
	return mediaController
}
