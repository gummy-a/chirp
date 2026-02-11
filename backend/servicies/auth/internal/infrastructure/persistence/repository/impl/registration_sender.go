package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/gummy_a/chirp/auth/internal/domain"
)

const (
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
	subject  = "Subject: Your register link is here!\n"
)

type RegistrationSenderRepository struct {
	logger *slog.Logger
}

func NewRegistrationSenderRepository(logger *slog.Logger) *RegistrationSenderRepository {
	return &RegistrationSenderRepository{logger: logger}
}

// ローカル確認用; gmailのapp passwordを事前に準備する必要あり
func (r *RegistrationSenderRepository) sendGmail(to_address domain.Email, token domain.NumberCode) error {
	from := os.Getenv("GMAIL_FROM_ADDRESS")
	if from == "" {
		return errors.New("GMAIL_FROM_ADDRESS is not set")
	}

	password := os.Getenv("GMAIL_APP_PASSWORD")
	if password == "" {
		return errors.New("GMAIL_APP_PASSWORD is not set")
	}

	body := "to register account, input following numbers " + fmt.Sprint(token) + "\n"
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to_address.String()}, message)
	if err != nil {
		r.logger.Error("Failed to send registration email", slog.String("to", to_address.String()), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *RegistrationSenderRepository) SendRegistrationEmail(to_address domain.Email, token domain.NumberCode) error {
	env := os.Getenv("APP_ENV")
	if env != "production" && env != "staging" {
		fmt.Printf("Registration token for %s: %d\n", to_address.String(), token)
		return nil
	}
	
	return r.sendGmail(to_address, token)
}
