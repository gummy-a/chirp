package controller

import (
	"context"
	"log/slog"

	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
)

func (h *AppHandler) ApiAuthV1LogoutPost(ctx context.Context, req api.ApiAuthV1LogoutPostRequest) (api.ImplResponse, error) {
	input := dto.ToLogoutInput(req)

	_, err := h.logoutUseCase.Execute(ctx, input)
	if err != nil {
		h.logger.Error("Failed to execute logout", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Unauthorized",
			Message: "failed to logout account",
		}}, nil
	}

	return api.ImplResponse{Code: 204, Body: nil}, nil
}
