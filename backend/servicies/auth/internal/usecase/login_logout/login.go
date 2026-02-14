package usecase

import (
	"context"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/usecase/repository"
)

type LoginAccountInput struct {
	Email    domain.Email
	Password domain.PasswordPlainText
}

type LoginAccountUseCase struct {
	account repository.AccountRepository
}

func NewLoginAccountUseCase(r repository.AccountRepository) *LoginAccountUseCase {
	return &LoginAccountUseCase{
		account: r,
	}
}

func (u *LoginAccountUseCase) Execute(ctx context.Context, input *LoginAccountInput) (*domain.JwtToken, error) {
	jwtToken, err := u.account.FindByEmailAndPassword(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
