package signupUsecase_test

import (
	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"chirp/backend/services/auth/v1/internal/usecase/_test/fake"
	usecase "chirp/backend/services/auth/v1/internal/usecase/signup"
	"testing"
	"time"
)

func TestSignupAccount_Success(t *testing.T) {
	tmpRepo := fake.NewTemporaryAccountRepo()
	accountRepo := fake.NewAccountRepo(tmpRepo)

	var id value_object.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := value_object.PasswordHash("password")
	numberCode := value_object.NumberCode(123456)
	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  value_object.Timestamp(time.Now().Add(1 * time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	token, err := uc.Execute(usecase.SignupAccountInput{
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

	var id value_object.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := value_object.PasswordHash("password")
	numberCode := value_object.NumberCode(123456)
	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  value_object.Timestamp(time.Now().Add(1 * time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	badNumberCode := value_object.NumberCode(987654)
	_, err = uc.Execute(usecase.SignupAccountInput{
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

	var id value_object.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := value_object.PasswordHash("password")
	numberCode := value_object.NumberCode(123456)
	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  value_object.Timestamp(time.Now().Add(-time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	_, err = uc.Execute(usecase.SignupAccountInput{
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

	var id value_object.TemporaryAccountID
	id.ParseString("6991c26a-8414-8324-9935-5b15cadb1c94")
	password := value_object.PasswordHash("password")
	numberCode := value_object.NumberCode(123456)
	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tmpAccount := &entity.TemporaryAccount{
		Id:         id,
		NumberCode: numberCode,
		ExpiresAt:  value_object.Timestamp(time.Now().Add(time.Hour)),
		Email:      email,
		Password:   password,
	}

	tmpRepo.Accounts[id.String()] = tmpAccount

	uc := usecase.NewSignupAccountUseCase(accountRepo, tmpRepo)

	_, err = uc.Execute(usecase.SignupAccountInput{
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
