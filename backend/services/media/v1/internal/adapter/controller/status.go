package mediaController

import (
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/usecase"
	"encoding/json"
	"net/http"
)

func (a *MediaController) Status(w http.ResponseWriter, r *http.Request) {
	jobId := r.PathValue("jobId")
	if jobId == "" {
		a.Logger.Error("jobId not found")
		http.Error(w, "jobId is required", http.StatusBadRequest)
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

	uuid, err := value_object.ParseUUID(jobId)
	if err != nil {
		a.Logger.Error("failed to parse jobId")
		http.Error(w, "invalid jobid", http.StatusBadRequest)
		return
	}
	id := value_object.JobId(uuid)

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

	err = a.UseCases.Status.Execute(&usecase.MediaEncodeStatusInput{
		ReqCtx:     r.Context(),
		JobId:      &id,
		WorkerFunc: workerFunc,
	})
	if err != nil {
		a.Logger.Error("failed to check status")
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}
}
