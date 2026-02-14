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
	Email    domain.Email
	Password domain.PasswordHash
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
	var hashedPassword domain.PasswordHash
	err := hashedPassword.NewHashFromBytes([]byte(input.Password))
	if err != nil {
		return nil, err
	}

	expiresAt := domain.Timestamp(time.Now().Add(24 * time.Hour))
	numberCode, tmpAccountID, err := u.repo.Create(ctx, input.Email, hashedPassword, expiresAt)
	if err != nil {
		return nil, err
	}

	err = u.reg.SendRegistrationEmail(&input.Email, numberCode, tmpAccountID)
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
