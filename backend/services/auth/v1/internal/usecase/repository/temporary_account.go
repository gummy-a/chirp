package repository

import (
	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
)

type TemporaryAccountRepository interface {
	Create(email value_object.Email, passwordHash value_object.PasswordHash, expiresAt value_object.Timestamp) (*value_object.NumberCode, *value_object.TemporaryAccountID, error)
	Delete(id value_object.TemporaryAccountID) error
	FindById(id value_object.TemporaryAccountID) (*entity.TemporaryAccount, error)
}
