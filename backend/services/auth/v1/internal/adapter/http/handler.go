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
			Path:    "POST /api/auth/v1/signup/",
			Handler: controller.PostSignup,
		},
		{
			Path:    "POST /api/auth/v1/signup/tmp/",
			Handler: controller.PostTmpSignup,
		},
		{
			Path:    "POST /api/auth/v1/login/",
			Handler: controller.PostLogin,
		},
		{
			Path:    "POST /api/auth/v1/account/tmp/{account_id}/",
			Handler: controller.GetTmpAccount,
		},
	}
	return authController
}
