package helper

import (
	"net/http"

	"github.com/gummy_a/chirp/media/internal/infrastructure/http/middleware"
)

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func NewChain(router http.Handler) http.Handler {

	/*
		router.Use() is often executed after the router has found the path.
		By wrapping the router itself, make sure CORS detection runs before the gorilla/mux router.
	*/
	middlewares := []func(http.Handler) http.Handler{
		middleware.JwtMiddleware,
		middleware.EnableCORS,
	}
	return chain(router, middlewares...)
}
