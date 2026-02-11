package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/auth/jwt"
)

type GenerateTokenRepository struct {
}

func NewGenerateTokenRepository() *GenerateTokenRepository {
	return &GenerateTokenRepository{}
}

func (r *GenerateTokenRepository) GenerateToken(accountID domain.AccountID) (*string, error) {
	token, err := jwt.GenerateJwt(accountID)
	if err != nil {
		return nil, err
	}
	return token, nil
}
