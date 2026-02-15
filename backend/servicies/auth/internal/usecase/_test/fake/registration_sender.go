package fake

import (
	"errors"
	
	"github.com/gummy_a/chirp/auth/internal/domain"
)

type RegistrationSender struct {
	Sent bool
}

func (f *RegistrationSender) SendRegistrationEmail(to_address *domain.Email, numberCode *domain.NumberCode, tmpAccountId *domain.TemporaryAccountID) error {
	f.Sent = true
	return nil
}

type FailingRegistrationSender struct{}

func (f *FailingRegistrationSender) SendRegistrationEmail(to_address *domain.Email, numberCode *domain.NumberCode, tmpAccountId *domain.TemporaryAccountID) error {
	return errors.New("send failed")
}