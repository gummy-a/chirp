package fake

import (
	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
)

type AccountRepo struct {
	Created bool
	TmpRepo *TemporaryAccountRepo
}

func NewAccountRepo(tmp *TemporaryAccountRepo) *AccountRepo {
	return &AccountRepo{
		TmpRepo: tmp,
	}
}

func (f *AccountRepo) CreateAccountThenDeleteTemporaryAccount(
	tmp entity.TemporaryAccount,
) (*value_object.JwtToken, error) {

	f.Created = true
	token := value_object.JwtToken("dummy.jwt.token")

	_ = f.TmpRepo.Delete(tmp.Id)
	return &token, nil
}

func (f *AccountRepo) Delete(id value_object.AccountID) error {
	return nil
}

func (f *AccountRepo) FindByEmailAndPassword(email value_object.Email, password value_object.PasswordPlainText) (*value_object.JwtToken, error) {
	token := value_object.JwtToken("dummy.jwt.token")
	return &token, nil
}

func (f *AccountRepo) FindFromJwtToken(jwtToken value_object.JwtToken) (*entity.Account, error) {

	var id value_object.AccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := value_object.PasswordHash("password")

	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		return nil, err
	}

	account := entity.Account{
		Id:       id,
		Email:    email,
		Password: password,
	}

	return &account, nil
}
