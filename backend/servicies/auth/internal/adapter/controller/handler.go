package controller

import (
	"log/slog"
	"net/http"

	"github.com/gummy_a/chirp/auth/internal/adapter/controller/helper"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
	loginLogoutUseCase "github.com/gummy_a/chirp/auth/internal/usecase/login_logout"
	signupUseCase "github.com/gummy_a/chirp/auth/internal/usecase/signup"
)

type AppHandler struct {
	loginUseCase     *loginLogoutUseCase.LoginAccountUseCase
	logoutUseCase    *loginLogoutUseCase.LogoutAccountUseCase
	tmpSignupUseCase *signupUseCase.SignupTemporaryAccountUseCase
	signupUseCase    *signupUseCase.SignupAccountUseCase
	logger           *slog.Logger
}

func NewAppRouter(tmpcase *signupUseCase.SignupTemporaryAccountUseCase, defcase *signupUseCase.SignupAccountUseCase, login *loginLogoutUseCase.LoginAccountUseCase, logout *loginLogoutUseCase.LogoutAccountUseCase, logger *slog.Logger) http.Handler {
	server := &AppHandler{
		tmpSignupUseCase: tmpcase,
		signupUseCase:    defcase,
		loginUseCase:     login,
		logoutUseCase:    logout,
		logger:           logger,
	}
	DefaultAPIController := api.NewDefaultAPIController(server)
	router := api.NewRouter(DefaultAPIController)

	return helper.NewChain(router)
}
