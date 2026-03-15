package loginUsecase

import (
	"chirp/backend/services/auth/v1/internal/domain/repository"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
)

type LoginAccountInput struct {
	Email    value_object.Email
	Password value_object.PasswordPlainText
}

type LoginAccountUseCase struct {
	account repository.AccountRepository
}

func NewLoginAccountUseCase(r repository.AccountRepository) *LoginAccountUseCase {
	return &LoginAccountUseCase{
		account: r,
	}
}

func (u *LoginAccountUseCase) Execute(input LoginAccountInput) (*value_object.JwtToken, error) {
	jwtToken, err := u.account.FindByEmailAndPassword(input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
