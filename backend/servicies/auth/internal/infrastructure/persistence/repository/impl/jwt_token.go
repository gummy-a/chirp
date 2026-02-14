package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/auth/jwt"
)

type TokenRepository struct {
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (r *TokenRepository) GenerateToken(accountID *domain.AccountID) (*string, error) {
	token, err := jwt.GenerateJwt(accountID)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *TokenRepository) ExtractClaims(token *domain.JwtToken) (*domain.AccountID, error) {
	accountId, err := jwt.ExtractClaims(token)
	if err != nil {
		return nil, err
	}
	return accountId, nil
}