package usecase_test

import (
	"context"
	"testing"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/usecase/_test/fake"
	"github.com/gummy_a/chirp/auth/internal/usecase/login_logout"
)

func TestLogoutTemporaryAccount_Success(t *testing.T) {
	tmp := fake.NewTemporaryAccountRepo()
	repo := fake.NewAccountRepo(tmp)

	uc := usecase.NewLogoutAccountUseCase(repo)

	jwtToken := domain.JwtToken("test")

	input := &usecase.LogoutAccountInput{
		JwtToken: jwtToken,
	}

	account, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if account == nil {
		t.Fatal("account should not be nil")
	}
}
