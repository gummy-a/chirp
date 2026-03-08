package entity

import "chirp/backend/services/auth/v1/internal/domain/value_object"

type Account struct {
	Id       value_object.AccountID
	Email    value_object.Email
	Password value_object.PasswordHash
}
