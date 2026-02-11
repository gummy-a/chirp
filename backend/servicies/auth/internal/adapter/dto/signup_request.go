package dto

import (
	api "github.com/gummy_a/chirp/auth/internal/adapter/openapi/signup/go"
	"github.com/gummy_a/chirp/auth/internal/domain"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/register_account"
)

func ToSignupTemporaryAccountInput(req api.ApiAuthV1TmpSignupPostRequest) *usecase.SignupTemporaryAccountInput {
	return &usecase.SignupTemporaryAccountInput{
		Email:    req.Email,
		Password: req.Password,
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
