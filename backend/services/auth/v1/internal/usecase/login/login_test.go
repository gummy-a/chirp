package loginUsecase_test

import (
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	usecase "chirp/backend/services/auth/v1/internal/usecase/login"
	"testing"

	"chirp/backend/services/auth/v1/internal/usecase/_test/fake"
)

func TestLoginTemporaryAccount_Success(t *testing.T) {
	tmp := fake.NewTemporaryAccountRepo()
	repo := fake.NewAccountRepo(tmp)

	uc := usecase.NewLoginAccountUseCase(repo)

	email, err := value_object.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	password := value_object.PasswordPlainText("password")

	input := usecase.LoginAccountInput{
		Email:    email,
		Password: password,
	}

	jwtToken, err := uc.Execute(input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if jwtToken == nil {
		t.Fatal("jwt token id should not be nil")
	}
}
