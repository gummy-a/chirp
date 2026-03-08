package entity

import "chirp/backend/services/auth/v1/internal/domain/value_object"

type TemporaryAccount struct {
	Id         value_object.TemporaryAccountID
	Email      value_object.Email
	Password   value_object.PasswordHash
	ExpiresAt  value_object.Timestamp
	NumberCode value_object.NumberCode
}
