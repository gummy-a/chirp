package authController

import (
	"chirp/backend/services/auth/v1/internal/usecase/login"
	"chirp/backend/services/auth/v1/internal/usecase/signup"
	"log/slog"
)

type AuthControllerUseCases struct {
	Signup    *signupUsecase.SignupAccountUseCase
	Login     *loginUsecase.LoginAccountUseCase
	TmpSignup *signupUsecase.SignupTemporaryAccountUseCase
}

type AuthController struct {
	UseCases *AuthControllerUseCases
	Logger   *slog.Logger
}
