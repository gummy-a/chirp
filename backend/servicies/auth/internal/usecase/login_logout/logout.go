package usecase

import (
	"context"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
	"github.com/gummy_a/chirp/auth/internal/usecase/repository"
)

type LogoutAccountInput struct {
	JwtToken domain.JwtToken
}

type LogoutAccountUseCase struct {
	account repository.AccountRepository
}

func NewLogoutAccountUseCase(r repository.AccountRepository) *LogoutAccountUseCase {
	return &LogoutAccountUseCase{
		account: r,
	}
}

func (u *LogoutAccountUseCase) Execute(ctx context.Context, input *LogoutAccountInput) (*entity.Account, error) {
	return u.account.FindFromJwtToken(ctx, &input.JwtToken)
}
