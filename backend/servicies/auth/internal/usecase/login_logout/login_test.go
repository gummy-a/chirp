package usecase_test

import (
	"context"
	"testing"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/usecase/_test/fake"
	"github.com/gummy_a/chirp/auth/internal/usecase/login_logout"
)

func TestLoginTemporaryAccount_Success(t *testing.T) {
	tmp := fake.NewTemporaryAccountRepo()
	repo := fake.NewAccountRepo(tmp)

	uc := usecase.NewLoginAccountUseCase(repo)

	email, err := domain.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := domain.PasswordPlainText("password")

	input := &usecase.LoginAccountInput{
		Email:    email,
		Password: password,
	}

	jwtToken, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if jwtToken == nil {
		t.Fatal("jwt token id should not be nil")
	}
}
