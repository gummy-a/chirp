package controller

import (
	"context"
	"log/slog"

	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
	loginLogoutUseCase "github.com/gummy_a/chirp/auth/internal/usecase/login_logout"
)

type LogoutHandler struct {
	logoutUseCase *loginLogoutUseCase.LogoutAccountUseCase
	logger        *slog.Logger
}

func NewLogoutHandler(logout *loginLogoutUseCase.LogoutAccountUseCase, logger *slog.Logger) *LogoutHandler {
	return &LogoutHandler{
		logoutUseCase: logout,
		logger:        logger,
	}
}

func (l *LogoutHandler) Logout(ctx context.Context, req api.ApiAuthV1LogoutPostRequest) (api.ImplResponse, error) {
	input := dto.ToLogoutInput(req)

	_, err := l.logoutUseCase.Execute(ctx, input)
	if err != nil {
		l.logger.Error("Failed to execute logout", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Unauthorized",
			Message: "failed to logout account",
		}}, nil
	}

	return api.ImplResponse{Code: 204, Body: nil}, nil
}
