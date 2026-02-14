package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type TokenRepository interface {
	GenerateToken(accountID *domain.AccountID) (*string, error)
	ExtractClaims(token *domain.JwtToken) (*domain.AccountID, error)
}
