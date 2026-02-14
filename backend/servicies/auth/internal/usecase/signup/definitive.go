package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/usecase/repository"
)

type SignupAccountInput struct {
	SignupToken domain.TemporaryAccountID
	NumberCode  domain.NumberCode
}

type SignupAccountUseCase struct {
	account    repository.AccountRepository
	tmpAccount repository.TemporaryAccountRepository
}

func NewSignupAccountUseCase(r repository.AccountRepository, tmp repository.TemporaryAccountRepository) *SignupAccountUseCase {
	return &SignupAccountUseCase{
		account:    r,
		tmpAccount: tmp,
	}
}

func (u *SignupAccountUseCase) Execute(ctx context.Context, input *SignupAccountInput) (*domain.JwtToken, error) {
	tempAccount, err := u.tmpAccount.FindById(ctx, &input.SignupToken)
	if err != nil {
		return nil, err
	}

	if tempAccount == nil {
		return nil, errors.New("temporary account not found")
	}

	if tempAccount.NumberCode != domain.NumberCode(input.NumberCode) {
		return nil, errors.New("invalid signup token")
	}

	if time.Now().After(time.Time(tempAccount.ExpiresAt)) {
		return nil, errors.New("signup token has expired")
	}

	jwtToken, err := u.account.CreateAccountThenDeleteTemporaryAccount(ctx, tempAccount)
	if err != nil {
		return nil, err
	}

	return jwtToken, nil
}
