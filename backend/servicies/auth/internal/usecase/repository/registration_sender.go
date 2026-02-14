package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type RegistrationSenderRepository interface {
	SendRegistrationEmail(to_address *domain.Email, numberCode *domain.NumberCode, tmpAccountId *domain.TemporaryAccountID) error
}
