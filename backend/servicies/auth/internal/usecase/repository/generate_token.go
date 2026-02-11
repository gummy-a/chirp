package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type GenerateTokenRepository interface {
	GenerateToken(accountID domain.AccountID) (*string, error)
	ValidateToken(token string) (*domain.AccountID, error)
}
