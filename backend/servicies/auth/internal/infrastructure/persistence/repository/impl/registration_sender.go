package repository

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gummy_a/chirp/auth/internal/domain"
)

type RegistrationSenderRepository struct {
	logger *slog.Logger
}

func NewRegistrationSenderRepository(logger *slog.Logger) *RegistrationSenderRepository {
	return &RegistrationSenderRepository{logger: logger}
}

func (r *RegistrationSenderRepository) SendRegistrationEmail(to_address *domain.Email, numberCode *domain.NumberCode, tmpAccountId *domain.TemporaryAccountID) error {
	// In non-production/staging environments, just log the token instead of sending an email
	env := os.Getenv("AUTH_SERVICE_APP_ENV")
	if env != "production" && env != "staging" {
		fmt.Printf("Registration token for %s: %d\n", to_address.String(), *numberCode)
		return nil
	}

	// TODO: use a legitimate email service in production
	return nil
}
