package handler

import (
	"chirp/backend/router"
	"chirp/backend/services/auth/v1/internal/adapter/controller"
	"log/slog"
)

func NewAuthHandler(uc *authController.AuthControllerUseCases, logger *slog.Logger) []router.Handler {
	controller := authController.AuthController{
		UseCases: uc,
		Logger:   logger,
	}
	authController := []router.Handler{
		{
			// PostSignup is a handler for definitive user registration.
			// Input: signup_token (string), number_code (int) from form values.
			// Success: 204 No Content, sets a session cookie.
			// Failure: 400 Bad Request for invalid token, number code, or signup failure.
			Path:    "POST /api/auth/v1/signup/",
			Handler: controller.PostSignup,
		},
		{
			// PostTmpSignup is a handler for temporary user registration.
			// Input: email (string), password (string) from form values.
			// Success: 200 OK with JSON body {"signup_token": "..."}.
			// Failure: 400 Bad Request for invalid input, 500 Internal Server Error for encoding errors.
			Path:    "POST /api/auth/v1/signup/tmp/",
			Handler: controller.PostTmpSignup,
		},
		{
			// PostLogin is a handler for user login.
			// Input: email (string), password (string) from form values.
			// Success: 204 No Content, sets a session cookie.
			// Failure: 400 Bad Request if login fails.
			Path:    "POST /api/auth/v1/login/",
			Handler: controller.PostLogin,
		},
		{
			// GetTmpAccount is a handler to get temporary account information.
			// Input: account_id from URL path.
			// Success: 204 No Content.
			// Failure: 400 Bad Request if account_id is missing, 404 Not Found if account not found, 500 Internal Server Error for parsing errors or other lookup failures.
			Path:    "GET /api/auth/v1/account/tmp/{account_id}/",
			Handler: controller.GetTmpAccount,
		},
	}
	return authController
}
