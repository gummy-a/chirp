package fake

import (
	"errors"

	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
)

type TemporaryAccountRepo struct {
	Accounts map[string]*entity.TemporaryAccount
}

func NewTemporaryAccountRepo() *TemporaryAccountRepo {
	return &TemporaryAccountRepo{
		Accounts: map[string]*entity.TemporaryAccount{},
	}
}

func (f *TemporaryAccountRepo) Create(
	email value_object.Email,
	password value_object.PasswordHash,
	expiresAt value_object.Timestamp,
) (*value_object.NumberCode, *value_object.TemporaryAccountID, error) {

	var id value_object.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")

	var numberCode value_object.NumberCode
	numberCode = value_object.NumberCode(123456)

	account := entity.TemporaryAccount{
		Id:         id,
		Email:      email,
		Password:   password,
		ExpiresAt:  expiresAt,
		NumberCode: numberCode,
	}

	f.Accounts[id.String()] = &account

	return &numberCode, &id, nil
}

func (f *TemporaryAccountRepo) FindById(
	id *value_object.TemporaryAccountID,
) (*entity.TemporaryAccount, error) {

	acc, ok := f.Accounts[id.String()]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}

func (f *TemporaryAccountRepo) Delete(
	id *value_object.TemporaryAccountID,
) error {

	delete(f.Accounts, id.String())
	return nil
}
