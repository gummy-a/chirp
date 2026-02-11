package entity

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type TemporaryAccount struct {
	Id         domain.TemporaryAccountID
	Email      domain.Email
	Password   domain.PasswordHash
	ExpiresAt  domain.Timestamp
	NumberCode domain.NumberCode
}
