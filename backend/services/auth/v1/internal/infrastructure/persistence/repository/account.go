package repository

import (
	"context"
	"log/slog"

	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"chirp/backend/services/auth/v1/internal/infrastructure/jwt"
	"chirp/backend/services/auth/v1/internal/infrastructure/persistence/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository struct {
	db     *pgxpool.Pool
	sql    *sqlc.Queries
	logger *slog.Logger
	ctx    *context.Context
}

func NewAccountRepository(db *pgxpool.Pool, q *sqlc.Queries, logger *slog.Logger, ctx *context.Context) *AccountRepository {
	return &AccountRepository{db: db, sql: q, logger: logger, ctx: ctx}
}

func (r *AccountRepository) CreateAccountThenDeleteTemporaryAccount(tmpAccount entity.TemporaryAccount) (*value_object.JwtToken, error) {
	email := value_object.Email(tmpAccount.Email)
	passwordHash := value_object.PasswordHash(tmpAccount.Password)
	algorithm := value_object.NewPasswordAlgorithm()
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(tmpAccount.Id),
		Valid: true,
	}

	// --- start a transaction
	transaction, err := r.db.Begin(*r.ctx)
	if err != nil {
		r.logger.Error("Failed to begin transaction", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}
	defer transaction.Rollback(*r.ctx)

	qtx := r.sql.WithTx(transaction)

	// create definitive account in the database
	createdAccount, err := qtx.CreateAccount(*r.ctx, sqlc.CreateAccountParams{
		Email:             email.String(),
		PasswordHash:      passwordHash.String(),
		PasswordAlgorithm: algorithm.String(),
	})
	if err != nil {
		r.logger.Error("Failed to create account", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}

	// delete temporary account
	_, err = qtx.DeleteTemporaryAccount(*r.ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to delete temporary account", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}

	// --- commit the transaction
	err = transaction.Commit(*r.ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}

	// generate JWT token for the new account
	account_id := value_object.AccountID(createdAccount.ID.Bytes)
	jwt, err := jwt.GenerateJwt(account_id)
	if err != nil {
		r.logger.Error("Failed to generate JWT", slog.String("account_id", account_id.String()), slog.String("error", err.Error()))
		return nil, err
	}

	jwtToken := value_object.JwtToken(*jwt)
	return &jwtToken, nil
}

func (r *AccountRepository) Delete(id value_object.AccountID) error {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}

	_, err := r.sql.DeleteAccount(*r.ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to delete account", slog.String("account_id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *AccountRepository) FindByEmailAndPassword(email value_object.Email, password value_object.PasswordPlainText) (*value_object.JwtToken, error) {
	account, err := r.sql.FindAccountByEmail(*r.ctx, email.String())
	if err != nil {
		r.logger.Error("Failed to find account by email and password", slog.String("email", account.Email), slog.String("password hash", account.PasswordHash), slog.String("error", err.Error()))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		r.logger.Error("Failed to compare password hash and password", slog.String("email", account.Email), slog.String("passwordHash", account.PasswordHash), slog.String("error", err.Error()))
		return nil, err
	}

	account_id := value_object.AccountID(account.ID.Bytes)
	jwt, err := jwt.GenerateJwt(account_id)
	if err != nil {
		r.logger.Error("Failed to generate JWT", slog.String("account_id", account_id.String()), slog.String("error", err.Error()))
		return nil, err
	}

	jwtToken := value_object.JwtToken(*jwt)
	return &jwtToken, nil

}
