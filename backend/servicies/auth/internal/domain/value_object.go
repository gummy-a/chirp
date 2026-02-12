package domain

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountID uuid.UUID
type TemporaryAccountID uuid.UUID
type Email string
type PasswordHash string
type PasswordAlgorithm string
type Timestamp time.Time
type JwtToken string
type NumberCode int32

func (id *AccountID) String() string {
	return uuid.UUID(*id).String()
}

func (id *TemporaryAccountID) String() string {
	return uuid.UUID(*id).String()
}

func NewTemporaryAccountIDFromSignupToken(s string) (TemporaryAccountID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return TemporaryAccountID(uuid.Nil), err
	}
	return TemporaryAccountID(parsed), nil
}

func (e *Email) String() string {
	return string(*e)
}

func NewEmail(v string) (Email, error) {
	_, err := mail.ParseAddress(v)
	if err != nil {
		return "", err
	}
	return Email(v), nil
}

func (p *PasswordHash) NewHashFromBytes(v []byte) error {
	hash, err := bcrypt.GenerateFromPassword(v, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*p = PasswordHash(string(hash))
	return nil
}

func (p *PasswordHash) String() string {
	return string(*p)
}

func NewPasswordAlgorithm() PasswordAlgorithm {
	// sha256, sha512, argon2, scrypt, bcrypt, etc. could be supported in the future
	return "Blowfish (bcrypt.GenerateFromPassword)"
}

func (p *PasswordAlgorithm) String() string {
	return string(*p)
}

func (j *JwtToken) String() string {
	return string(*j)
}
