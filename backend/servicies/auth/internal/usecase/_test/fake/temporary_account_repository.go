package fake

import (
	"context"
	"errors"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
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
	ctx context.Context,
	email domain.Email,
	password domain.PasswordHash,
	expiresAt domain.Timestamp,
) (*domain.NumberCode, *domain.TemporaryAccountID, error) {

	var id domain.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")

	var numberCode domain.NumberCode
	numberCode = domain.NumberCode(123456)

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
	ctx context.Context,
	id *domain.TemporaryAccountID,
) (*entity.TemporaryAccount, error) {

	acc, ok := f.Accounts[id.String()]
	if !ok {
		return nil, errors.New("not found")
	}
	return acc, nil
}

func (f *TemporaryAccountRepo) Delete(
	ctx context.Context,
	id *domain.TemporaryAccountID,
) error {

	delete(f.Accounts, id.String())
	return nil
}