package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type RegistrationSenderRepository interface {
	SendRegistrationEmail(to_address domain.Email, token domain.NumberCode) error
}
