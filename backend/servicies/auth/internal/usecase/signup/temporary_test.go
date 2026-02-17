package usecase_test

import (
	"context"
	"testing"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/usecase/_test/fake"
	"github.com/gummy_a/chirp/auth/internal/usecase/signup"
)

func TestSignupTemporaryAccount_Success(t *testing.T) {
	// Arrange
	repo := fake.NewTemporaryAccountRepo()
	sender := &fake.RegistrationSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := domain.PasswordHash("password")

	input := &usecase.SignupTemporaryAccountInput{
		Email:    email,
		Password: password,
	}

	// test Execute
	id, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id == nil {
		t.Fatal("temporary account id should not be nil")
	}

	if !sender.Sent {
		t.Fatal("registration email should be sent")
	}

	if _, ok := repo.Accounts[id.String()]; !ok {
		t.Fatal("temporary account should be saved in repository")
	}
}

func TestFindByIdTemporaryAccount_Success(t *testing.T) {
	repo := fake.NewTemporaryAccountRepo()
	sender := &fake.RegistrationSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := domain.PasswordHash("password")

	input := &usecase.SignupTemporaryAccountInput{
		Email:    email,
		Password: password,
	}

	accountId, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if accountId == nil {
		t.Fatal("temporary account should not be found")
	}

	account, err := uc.FindById(context.Background(), accountId)

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
	sender := &fake.RegistrationSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	var id domain.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")

	account, err := uc.FindById(context.Background(), &id)

	if err == nil {
		t.Fatal("temporary account should not be found")
	}

	if account != nil {
		t.Fatal("temporary account should be nil")
	}
}

func TestSignupTemporaryAccount_EmailFailure_DeletesAccount(t *testing.T) {
	repo := fake.NewTemporaryAccountRepo()
	sender := &fake.FailingRegistrationSender{}

	uc := usecase.NewSignupTemporaryAccountUseCase(repo, sender)

	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := domain.PasswordHash("password")

	input := &usecase.SignupTemporaryAccountInput{
		Email:    email,
		Password: password,
	}
	_, err = uc.Execute(context.Background(), input)

	if err == nil {
		t.Fatal("expected error")
	}

	if len(repo.Accounts) != 0 {
		t.Fatal("temporary account should be deleted on failure")
	}
}
