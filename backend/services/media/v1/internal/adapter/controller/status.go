package mediaController

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/infrastructure/http/middleware"
	"chirp/backend/services/media/v1/internal/usecase"
	"encoding/json"
	"net/http"
)

func (a *MediaController) Status(w http.ResponseWriter, r *http.Request) {
	ownerAccountId, ok := r.Context().Value(middleware.OwnerAccountIdKey).(value_object.OwnerAccountId)
	if !ok {
		a.Logger.Error("Failed to get accountId from context")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		a.Logger.Error("failed to start SSE")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	workerFunc := func(data map[string]interface{}) {
		w.Write([]byte("data: ")) // SSE must start with data: ...
		b, err := json.Marshal(data)
		if err != nil {
			a.Logger.Error("failed to marshal data")
			w.Write([]byte("internal error"))
		} else {
			w.Write(b)
		}
		w.Write([]byte("\n\n")) // SSE must end with \n\n
		flusher.Flush()
	}

	err := a.UseCases.Status.Execute(&usecase.MediaEncodeStatusInput{
		ReqCtx:         r.Context(),
		OwnerAccountId: &ownerAccountId,
		WorkerFunc:     workerFunc,
	})
	if err != nil {
		a.Logger.Error("failed to check status")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
}
