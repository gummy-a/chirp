package router

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gummy_a/chirp/auth/internal/adapter/controller"
	"github.com/gummy_a/chirp/auth/internal/adapter/controller/helper"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
)

type AppHandler struct {
	loginHandler  *controller.LoginHandler
	signupHandler *controller.SignupHandler
	logoutHandler *controller.LogoutHandler
	logger        *slog.Logger
}

func NewAppRouter(login *controller.LoginHandler, logout *controller.LogoutHandler, signup *controller.SignupHandler, logger *slog.Logger) http.Handler {
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

func (h *AppHandler) ApiAuthV1LogoutPost(ctx context.Context, req api.ApiAuthV1LogoutPostRequest) (api.ImplResponse, error) {
	return h.logoutHandler.Logout(ctx, req)
}

func (h *AppHandler) ApiAuthV1LoginPost(ctx context.Context, req api.ApiAuthV1LoginPostRequest) (api.ImplResponse, error) {
	return h.loginHandler.Login(ctx, req)
}

func (h *AppHandler) ApiAuthV1TmpSignupPost(ctx context.Context, req api.ApiAuthV1TmpSignupPostRequest) (api.ImplResponse, error) {
	return h.signupHandler.TmpSignup(ctx, req)
}

func (h *AppHandler) ApiAuthV1SignupPost(ctx context.Context, req api.ApiAuthV1SignupPostRequest) (api.ImplResponse, error) {
	return h.signupHandler.Signup(ctx, req)
}

func (h *AppHandler) ApiAuthV1TmpAccountIdGet(ctx context.Context, id string) (api.ImplResponse, error) {
	return h.signupHandler.FindTemporaryAccountById(ctx, id)
}
