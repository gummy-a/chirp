package signupUsecase_test

import (
	"testing"

	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"chirp/backend/services/auth/v1/internal/usecase/_test/fake"
	usecase "chirp/backend/services/auth/v1/internal/usecase/signup"
)

type fakeEmailSender struct {
	err error
}

func (f *fakeEmailSender) Send(to value_object.Email, numberCode value_object.NumberCode, tmpAccountID *value_object.TemporaryAccountID) error {
	return f.err
}

func TestSignupTemporaryAccount_Success(t *testing.T) {
	// Arrange
	repo := fake.NewTemporaryAccountRepo()
	sender := &fakeEmailSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := value_object.PasswordHash("password")

	input := usecase.SignupTemporaryAccountInput{
		Email:    email,
		Password: password,
	}

	// test Execute
	id, err := uc.Execute(&input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id == nil {
		t.Fatal("temporary account id should not be nil")
	}

	if _, ok := repo.Accounts[id.String()]; !ok {
		t.Fatal("temporary account should be saved in repository")
	}
}

func TestFindByIdTemporaryAccount_Success(t *testing.T) {
	repo := fake.NewTemporaryAccountRepo()
	sender := &fakeEmailSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := value_object.PasswordHash("password")

	input := usecase.SignupTemporaryAccountInput{
		Email:    email,
		Password: password,
	}

	accountId, err := uc.Execute(&input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if accountId == nil {
		t.Fatal("temporary account should not be found")
	}

	account, err := uc.FindById(accountId)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if account == nil {
		t.Fatal("temporary account should not be nil")
	}

	if account.Email != email {
		t.Fatal("FindById should find the created one")
	}
}

func TestFindByIdTemporaryAccount_FailureBeforeSignup(t *testing.T) {
	repo := fake.NewTemporaryAccountRepo()
	sender := &fakeEmailSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	var id value_object.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")

	account, err := uc.FindById(&id)

	if err == nil {
		t.Fatal("temporary account should not be found")
	}

	if account != nil {
		t.Fatal("temporary account should be nil")
	}
}
