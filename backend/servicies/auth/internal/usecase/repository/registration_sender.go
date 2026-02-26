package repository

import (
	"github.com/gummy_a/chirp/auth/internal/domain"
)

// TODO: refactor me;
// Email sending is not appropriate for a repository
type RegistrationSenderRepository interface {
	SendRegistrationEmail(to_address *domain.Email, numberCode *domain.NumberCode, tmpAccountId *domain.TemporaryAccountID) error
}
