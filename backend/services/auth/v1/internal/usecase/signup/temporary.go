package signupUsecase

import (
	"time"

	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/repository"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"chirp/backend/services/auth/v1/internal/usecase/port/email"
)

type SignupTemporaryAccountInput struct {
	Email    value_object.Email
	Password value_object.PasswordHash
}

type SignupTemporaryAccountUseCase struct {
	repo  repository.TemporaryAccountRepository
	email email.EmailSender
}

func NewSignupTemporaryAccountUseCase(r repository.TemporaryAccountRepository, e email.EmailSender) *SignupTemporaryAccountUseCase {
	return &SignupTemporaryAccountUseCase{
		repo:  r,
		email: e,
	}
}

func (u *SignupTemporaryAccountUseCase) Execute(input *SignupTemporaryAccountInput) (*value_object.TemporaryAccountID, error) {
	var hashedPassword value_object.PasswordHash
	err := hashedPassword.NewHashFromBytes([]byte(input.Password))
	if err != nil {
		return nil, err
	}

	expiresAt := value_object.Timestamp(time.Now().Add(24 * time.Hour))
	numberCode, tmpAccountID, err := u.repo.Create(input.Email, hashedPassword, expiresAt)
	if err != nil {
		return nil, err
	}

	err = u.email.Send(input.Email, *numberCode, tmpAccountID)
	if err != nil {
		return nil, err
	}

	return tmpAccountID, nil
}

func (u *SignupTemporaryAccountUseCase) FindById(id *value_object.TemporaryAccountID) (*entity.TemporaryAccount, error) {
	return u.repo.FindById(id)
}
