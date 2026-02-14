package controller

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	"github.com/gummy_a/chirp/auth/internal/adapter/middleware"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
)

func (h *AppHandler) ApiAuthV1LoginPost(ctx context.Context, req api.ApiAuthV1LoginPostRequest) (api.ImplResponse, error) {
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
