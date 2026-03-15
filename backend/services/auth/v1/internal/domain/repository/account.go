package repository

import (
	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
)

type AccountRepository interface {
	CreateAccountThenDeleteTemporaryAccount(account *entity.TemporaryAccount) (*value_object.JwtToken, error)
	Delete(id *value_object.AccountID) error
	FindByEmailAndPassword(email value_object.Email, password value_object.PasswordPlainText) (*value_object.JwtToken, error)
}
