package router

import (
	"log"
	"net/http"
	"time"
)

type Handler struct {
	Path    string
	Handler http.HandlerFunc
}

func NewAppRouter(handler []Handler) *http.ServeMux {
	mux := http.NewServeMux()

	for _, h := range handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s \n", r.Method, r.RequestURI)
			s := time.Now()

			h.Handler(w, r)

			log.Printf("%s %s elapsed: %v\n", r.Method, r.RequestURI, time.Since(s))
		}
		mux.HandleFunc(h.Path, fn)
	}

	return mux
}
