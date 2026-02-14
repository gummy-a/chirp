package dto

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
	api "github.com/gummy_a/chirp/auth/internal/infrastructure/auth/openapi/auth/go"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/signup"
)

func ToSignupTemporaryAccountInput(req api.ApiAuthV1TmpSignupPostRequest) *usecase.SignupTemporaryAccountInput {
	return &usecase.SignupTemporaryAccountInput{
		Email:    domain.Email(req.Email),
		Password: domain.PasswordHash(req.Password),
	}
}

func ToSignupAccountInput(req api.ApiAuthV1SignupPostRequest) (*usecase.SignupAccountInput, error) {
	tmpAccountID, err := domain.NewTemporaryAccountIDFromSignupToken(req.SignupToken)
	if err != nil {
		return nil, err
	}

	return &usecase.SignupAccountInput{
		SignupToken: tmpAccountID,
		NumberCode:  domain.NumberCode(req.NumberCode),
	}, nil
}

func ToTokenFromAccountId(id domain.TemporaryAccountID) string {
	return id.String()
}
