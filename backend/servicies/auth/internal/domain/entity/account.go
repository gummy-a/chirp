package entity

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type Account struct {
	Id       domain.AccountID
	Email    domain.Email
	Password domain.PasswordHash
}
