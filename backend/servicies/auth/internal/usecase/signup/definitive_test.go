package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
	"github.com/gummy_a/chirp/auth/internal/usecase/_test/fake"
	"github.com/gummy_a/chirp/auth/internal/usecase/signup"
)

func TestSignupAccount_Success(t *testing.T) {
	tmpRepo := fake.NewTemporaryAccountRepo()
	accountRepo := fake.NewAccountRepo(tmpRepo)

	var id domain.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := domain.PasswordHash("password")
	numberCode := domain.NumberCode(123456)
	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  domain.Timestamp(time.Now().Add(1 * time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	token, err := uc.Execute(context.Background(), &usecase.SignupAccountInput{
		SignupToken: id,
		NumberCode:  numberCode,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == nil {
		t.Fatal("jwt token should be returned")
	}

	if !accountRepo.Created {
		t.Fatal("account should be created")
	}
}

func TestSignupAccount_InvalidNumberCode(t *testing.T) {
	tmpRepo := fake.NewTemporaryAccountRepo()
	accountRepo := fake.NewAccountRepo(tmpRepo)

	var id domain.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := domain.PasswordHash("password")
	numberCode := domain.NumberCode(123456)
	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  domain.Timestamp(time.Now().Add(1 * time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	badNumberCode := domain.NumberCode(987654)
	_, err = uc.Execute(context.Background(), &usecase.SignupAccountInput{
		SignupToken: id,
		NumberCode:  badNumberCode,
	})

	if err == nil {
		t.Fatal("expected error for invalid number code")
	}

	if accountRepo.Created {
		t.Fatal("account should not be created")
	}
}

func TestSignupAccount_ExpiredToken(t *testing.T) {
	tmpRepo := fake.NewTemporaryAccountRepo()
	accountRepo := fake.NewAccountRepo(tmpRepo)

	var id domain.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := domain.PasswordHash("password")
	numberCode := domain.NumberCode(123456)
	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  domain.Timestamp(time.Now().Add(-time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	_, err = uc.Execute(context.Background(), &usecase.SignupAccountInput{
		SignupToken: id,
		NumberCode:  numberCode,
	})

	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestSignupAccount_DeletesTemporaryAccount(t *testing.T) {
	tmpRepo := fake.NewTemporaryAccountRepo()
	accountRepo := fake.NewAccountRepo(tmpRepo)

	var id domain.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := domain.PasswordHash("password")
	numberCode := domain.NumberCode(123456)
	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  domain.Timestamp(time.Now().Add(time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	_, err = uc.Execute(context.Background(), &usecase.SignupAccountInput{
		SignupToken: id,
		NumberCode:  numberCode,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := tmpRepo.Accounts[id.String()]; exists {
		t.Fatal("temporary account should be deleted after successful signup")
	}
}
