package usecase

import (
	"context"
	"log"
	"time"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
	"github.com/gummy_a/chirp/auth/internal/usecase/repository"
)

type SignupTemporaryAccountInput struct {
	Email    string
	Password string
}

type SignupTemporaryAccountUseCase struct {
	repo repository.TemporaryAccountRepository
	reg  repository.RegistrationSenderRepository
}

func NewSignupTemporaryAccountUseCase(r repository.TemporaryAccountRepository, reg repository.RegistrationSenderRepository) *SignupTemporaryAccountUseCase {
	return &SignupTemporaryAccountUseCase{
		repo: r,
		reg:  reg,
	}
}

func (u *SignupTemporaryAccountUseCase) Execute(ctx context.Context, input *SignupTemporaryAccountInput) (*domain.TemporaryAccountID, error) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	var hashedPassword domain.PasswordHash
	err = hashedPassword.NewHashFromBytes([]byte(input.Password))
	if err != nil {
		return nil, err
	}

	expiresAt := domain.Timestamp(time.Now().Add(24 * time.Hour))
	token, tmpAccountID, err := u.repo.Create(ctx, email, hashedPassword, expiresAt)
	if err != nil {
		return nil, err
	}

	err = u.reg.SendRegistrationEmail(email, *token)
	if err != nil {
		inner_error := u.repo.Delete(ctx, tmpAccountID)
		if inner_error != nil {
			log.Printf("failed to delete temporary account after email send failure: %v", inner_error)
		}
		return nil, err
	}

	return tmpAccountID, nil
}

func (u *SignupTemporaryAccountUseCase) FindById(ctx context.Context, id *domain.TemporaryAccountID) (*entity.TemporaryAccount, error) {
	return u.repo.FindById(ctx, id)
}
