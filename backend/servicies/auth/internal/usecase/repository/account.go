package repository

import (
	"context"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
)

type AccountRepository interface {
	CreateAccountThenDeleteTemporaryAccount(ctx context.Context, account *entity.TemporaryAccount) (*domain.JwtToken, error)
	Delete(ctx context.Context, id domain.AccountID) error
	FindByEmailAndPassword(ctx context.Context, email domain.Email, password domain.PasswordPlainText) (*domain.JwtToken, error)
	FindFromJwtToken(ctx context.Context, jwtToken *domain.JwtToken) (*entity.Account, error)
}
