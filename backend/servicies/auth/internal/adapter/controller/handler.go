package controller

import (
	"log/slog"
	"net/http"

	"github.com/gummy_a/chirp/auth/internal/adapter/controller/helper"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
)

type AppHandler struct {
	loginHandler  *LoginHandler
	signupHandler *SignupHandler
	logoutHandler *LogoutHandler
	logger        *slog.Logger
}

func NewAppRouter(login *LoginHandler, logout *LogoutHandler, signup *SignupHandler, logger *slog.Logger) http.Handler {
	server := &AppHandler{
		signupHandler: signup,
		loginHandler:  login,
		logoutHandler: logout,
		logger:        logger,
	}
	DefaultAPIController := api.NewDefaultAPIController(server)
	router := api.NewRouter(DefaultAPIController)

	return helper.NewChain(router)
}
