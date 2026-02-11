package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	"github.com/gummy_a/chirp/auth/internal/adapter/middleware"
	api "github.com/gummy_a/chirp/auth/internal/adapter/openapi/signup/go"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/register_account"
)


func NewSignupRouter(tmpcase *usecase.SignupTemporaryAccountUseCase, defcase *usecase.SignupAccountUseCase, logger *slog.Logger) *mux.Router {
	server := &signupHandler{
		tmpSignupUseCase: tmpcase,
		signupUseCase: defcase,
		logger: logger,
	}
	DefaultAPIController := api.NewDefaultAPIController(server)
	router := api.NewRouter(DefaultAPIController)
    router.Use(middleware.MiddlewareStoreWriter)
	return router
}

type signupHandler struct {
	tmpSignupUseCase *usecase.SignupTemporaryAccountUseCase
	signupUseCase *usecase.SignupAccountUseCase
	logger *slog.Logger
}

func (h *signupHandler) AuthV1TmpSignupPost(ctx context.Context, req api.AuthV1TmpSignupPostRequest) (api.ImplResponse, error) {
	accountId, err := h.tmpSignupUseCase.Execute(ctx, dto.ToSignupTemporaryAccountInput(req))
	if err != nil {
		h.logger.Error("Failed to execute temporary account signup", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to register temporary account",
		}}, err
	}

	return api.ImplResponse{Code: 200, Body: api.AuthV1TmpSignupPost200Response{
		SignupToken: dto.ToTokenFromAccountId(*accountId),
	}}, nil
}

func (h *signupHandler) AuthV1SignupPost(ctx context.Context, req api.AuthV1SignupPostRequest) (api.ImplResponse, error) {
	ToSignupAccountInput, err := dto.ToSignupAccountInput(req)
	if err != nil {
		h.logger.Error("Failed to convert signup account input", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid signup token",
		}}, err
	}

	token, err := h.signupUseCase.Execute(ctx, ToSignupAccountInput)
	if err != nil {
		h.logger.Error("Failed to execute account signup", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 401, Body: api.ErrorResponse{
			Error:   "Unauthorized",
			Message: "failed to register account",
		}}, err
	}

    if w, ok := ctx.Value(middleware.ResponseWriterKey).(http.ResponseWriter); ok {
        http.SetCookie(w, &http.Cookie{
            Name:     "token",
            Value:    token.String(),
            Path:     "/",
            HttpOnly: true,
            Secure:   os.Getenv("APP_ENV") == "production",
            SameSite: http.SameSiteLaxMode,
            MaxAge:   86400,          // 24h
        })
    } else {
		h.logger.Error("ResponseWriter not found in context")
	}

	return api.ImplResponse{Code: 204, Body: nil}, nil
}
