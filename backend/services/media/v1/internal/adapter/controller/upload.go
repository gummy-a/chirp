package mediaController

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/infrastructure/http/middleware"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/file"
	"chirp/backend/services/media/v1/internal/usecase"
)

type UploadResponse struct {
	OriginalFileName string `json:"original_file_name"`
	JobId            string `json:"job_id"`
}

func (a *MediaController) Upload(w http.ResponseWriter, r *http.Request) {
	ownerAccountId, ok := r.Context().Value(middleware.OwnerAccountIdKey).(value_object.OwnerAccountId)
	if !ok {
		a.Logger.Error("Failed to get accountId from context")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		a.Logger.Error("Failed to ParseMultipartForm()", slog.String("error", err.Error()))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(w, "Failed to upload files", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	uploadedFile, err := file.SaveFilesToTmp(files)
	if err != nil {
		a.Logger.Error("Failed to SaveFilesToTmp()", slog.String("error", err.Error()))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	output, err := a.UseCases.Upload.EnqueueEncode(&usecase.MediaUploadInput{
		Files:          *uploadedFile,
		OwnerAccountId: ownerAccountId,
	})
	if err != nil {
		a.Logger.Error("Failed to EnqueueEncode", slog.String("error", err.Error()))
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return
	}

	var response []UploadResponse
	for _, v := range *output {
		response = append(response, UploadResponse{
			OriginalFileName: string(v.OriginalFileName),
			JobId:            v.JobId.String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		a.Logger.Error("Failed to encode response", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
