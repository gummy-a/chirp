package email

import "chirp/backend/services/auth/v1/internal/domain/value_object"

type EmailSender interface {
	Send(to value_object.Email, numberCode value_object.NumberCode, tmpAccountID *value_object.TemporaryAccountID) error
}
