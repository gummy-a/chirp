package controller

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	"github.com/gummy_a/chirp/auth/internal/adapter/middleware"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
	loginLogoutUseCase "github.com/gummy_a/chirp/auth/internal/usecase/login_logout"
)

type LoginHandler struct {
	loginUseCase *loginLogoutUseCase.LoginAccountUseCase
	logger       *slog.Logger
}

func NewLoginHandler(Login *loginLogoutUseCase.LoginAccountUseCase, logger *slog.Logger) *LoginHandler {
	return &LoginHandler{
		loginUseCase: Login,
		logger:       logger,
	}
}

func (h *LoginHandler) Login(ctx context.Context, req api.ApiAuthV1LoginPostRequest) (api.ImplResponse, error) {
	input := dto.ToLoginInput(req)
	jwtToken, err := h.loginUseCase.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Failed to execute login", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Unauthorized",
			Message: "failed to login account",
		}}, nil
	}

	if rw, ok := ctx.Value(middleware.ResponseWriterKey).(http.ResponseWriter); ok {
		cookie := &http.Cookie{
			Name:     "session",
			Value:    jwtToken.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   os.Getenv("AUTH_SERVICE_APP_ENV") == "production",
			MaxAge:   60 * 60 * 24, // 1day
		}
		http.SetCookie(rw, cookie)
	}

	return api.ImplResponse{Code: 204, Body: nil}, nil
}
