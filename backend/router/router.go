package router

import "net/http"

type Handler struct {
	Path    string
	Handler http.HandlerFunc
}

func NewAppRouter(handler []Handler) *http.ServeMux {
	mux := http.NewServeMux()

	for _, h := range handler {
		mux.HandleFunc(h.Path, h.Handler)
	}

	return mux
}
