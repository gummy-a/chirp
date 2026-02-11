package repository

import (
	"context"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
)

type TemporaryAccountRepository interface {
	Create(ctx context.Context, email domain.Email, passwordHash domain.PasswordHash, expiresAt domain.Timestamp) (*domain.NumberCode, *domain.TemporaryAccountID, error)
	Delete(ctx context.Context, id *domain.TemporaryAccountID) error
	FindById(ctx context.Context, id *domain.TemporaryAccountID) (*entity.TemporaryAccount, error)
}
