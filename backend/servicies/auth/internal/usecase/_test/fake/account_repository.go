package fake

import (
	"context"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
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
	ctx context.Context,
	tmp *entity.TemporaryAccount,
) (*domain.JwtToken, error) {

	f.Created = true
	token := domain.JwtToken("dummy.jwt.token")

	_ = f.TmpRepo.Delete(ctx, &tmp.Id)
	return &token, nil
}

func (f *AccountRepo)Delete(ctx context.Context, id domain.AccountID) error {
	return nil
}

func (f *AccountRepo)FindByEmailAndPassword(ctx context.Context, email domain.Email, password domain.PasswordPlainText) (*domain.JwtToken, error) {
	token := domain.JwtToken("dummy.jwt.token")
	return &token, nil
}

func (f *AccountRepo)FindFromJwtToken(ctx context.Context, jwtToken *domain.JwtToken) (*entity.Account, error) {

	var id domain.AccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := domain.PasswordHash("password")

	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		return nil, err
	}

	account := entity.Account{
		Id:         id,
		Email:      email,
		Password:   password,
	}

	return &account, nil
}
