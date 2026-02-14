package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gummy_a/chirp/auth/internal/adapter/dto"
	"github.com/gummy_a/chirp/auth/internal/adapter/middleware"
	api "github.com/gummy_a/chirp/auth/internal/adapter/openapi/signup/go"
	"github.com/gummy_a/chirp/auth/internal/domain"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/register_account"
)

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func NewSignupRouter(tmpcase *usecase.SignupTemporaryAccountUseCase, defcase *usecase.SignupAccountUseCase, logger *slog.Logger) http.Handler {
	server := &signupHandler{
		tmpSignupUseCase: tmpcase,
		signupUseCase:    defcase,
		logger:           logger,
	}
	DefaultAPIController := api.NewDefaultAPIController(server)
	router := api.NewRouter(DefaultAPIController)

	/*
		router.Use() is often executed after the router has found the path.
		By wrapping the router itself, make sure CORS detection runs before the gorilla/mux router.
	*/
	middlewares := []func(http.Handler) http.Handler{
		middleware.MiddlewareStoreWriter,
		middleware.EnableCORS,
	}
	return chain(router, middlewares...)
}

type signupHandler struct {
	tmpSignupUseCase *usecase.SignupTemporaryAccountUseCase
	signupUseCase    *usecase.SignupAccountUseCase
	logger           *slog.Logger
}

func (h *signupHandler) ApiAuthV1TmpSignupPost(ctx context.Context, req api.ApiAuthV1TmpSignupPostRequest) (api.ImplResponse, error) {
	accountId, err := h.tmpSignupUseCase.Execute(ctx, dto.ToSignupTemporaryAccountInput(req))
	if err != nil {
		h.logger.Error("Failed to execute temporary account signup", slog.String("error", err.Error()))
		return api.ImplResponse{Code: 400, Body: api.ErrorResponse{
			Error:   "Bad Request",
			Message: "failed to register temporary account",
		}}, nil // set return value 'error' nil to avoid returning plain text error message
	}

	return api.ImplResponse{Code: 200, Body: api.ApiAuthV1TmpSignupPost200Response{
		SignupToken: dto.ToTokenFromAccountId(*accountId),
	}}, nil
}

func (h *signupHandler) ApiAuthV1SignupPost(ctx context.Context, req api.ApiAuthV1SignupPostRequest) (api.ImplResponse, error) {
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
			Message: "failed to register account",
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

func (h *signupHandler) ApiAuthV1TmpAccountIdGet(ctx context.Context, id string) (api.ImplResponse, error) {
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
