package controller

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	"github.com/gummy_a/chirp/auth/internal/adapter/middleware"
	"github.com/gummy_a/chirp/auth/internal/domain"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
)

func (h *AppHandler) ApiAuthV1TmpSignupPost(ctx context.Context, req api.ApiAuthV1TmpSignupPostRequest) (api.ImplResponse, error) {
	accountId, err := h.tmpSignupUseCase.Execute(ctx, dto.ToSignupTemporaryAccountInput(req))
	if err != nil {
		h.logger.Error("Failed to execute temporary account signup", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to signup temporary account",
		}}, nil // set return value 'error' nil to avoid returning plain text error message
	}

	return api.ImplResponse{Code: 200, Body: api.ApiAuthV1TmpSignupPost200Response{
		SignupToken: dto.ToTokenFromAccountId(*accountId),
	}}, nil
}

func (h *AppHandler) ApiAuthV1SignupPost(ctx context.Context, req api.ApiAuthV1SignupPostRequest) (api.ImplResponse, error) {
	ToSignupAccountInput, err := dto.ToSignupAccountInput(req)
	if err != nil {
		h.logger.Error("Failed to convert signup account input", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid signup token",
		}}, nil
	}

	token, err := h.signupUseCase.Execute(ctx, ToSignupAccountInput)
	if err != nil {
		h.logger.Error("Failed to execute account signup", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Unauthorized",
			Message: "failed to signup account",
		}}, nil
	}

	if rw, ok := ctx.Value(middleware.ResponseWriterKey).(http.ResponseWriter); ok {
		cookie := &http.Cookie{
			Name:     "session",
			Value:    token.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   os.Getenv("AUTH_SERVICE_APP_ENV") == "production",
			MaxAge:   60 * 60 * 24, // 1day
		}
		http.SetCookie(rw, cookie)
	}

	return api.ImplResponse{Code: 204, Body: nil}, nil
}

func (h *AppHandler) ApiAuthV1TmpAccountIdGet(ctx context.Context, id string) (api.ImplResponse, error) {
	var domainId domain.TemporaryAccountID
	err := domainId.ParseString(id)
	if err != nil {
		h.logger.Error("Failed to parse temporary account id", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 404, Body: nil}, nil
	}

	accountId, err := h.tmpSignupUseCase.FindById(ctx, &domainId)
	if err != nil {
		h.logger.Error("Failed to get temporary account id", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 404, Body: nil}, nil
	}

	if accountId == nil {
		return api.ImplResponse{Code: 404, Body: nil}, nil
	}

	return api.ImplResponse{Code: 204, Body: nil}, nil
}
