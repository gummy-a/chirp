package dto

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/login_logout"
)

func ToLoginInput(req api.ApiAuthV1LoginPostRequest) *usecase.LoginAccountInput {
	return &usecase.LoginAccountInput{
		Email:    domain.Email(req.Email),
		Password: domain.PasswordPlainText(req.Password),
	}
}

func ToLogoutInput(req api.ApiAuthV1LogoutPostRequest) (*usecase.LogoutAccountInput) {
	return &usecase.LogoutAccountInput{
		JwtToken: domain.JwtToken(req.Session),
	}
}
