package repository

import (
	"context"
	"log/slog"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/auth/jwt"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db     *pgxpool.Pool
	sql    *sqlc.Queries
	logger *slog.Logger
}

func NewAccountRepository(db *pgxpool.Pool, q *sqlc.Queries, logger *slog.Logger) *AccountRepository {
	return &AccountRepository{db: db, sql: q, logger: logger}
}

func (r *AccountRepository) CreateAccountThenDeleteTemporaryAccount(ctx context.Context, tmpAccount *entity.TemporaryAccount) (*domain.JwtToken, error) {
	email := domain.Email(tmpAccount.Email)
	passwordHash := domain.PasswordHash(tmpAccount.Password)
	algorithm := domain.NewPasswordAlgorithm()

	// --- start a transaction
	transaction, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Error("Failed to begin transaction", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}
	defer transaction.Rollback(ctx)

	qtx := r.sql.WithTx(transaction)

	// create definitive account in the database
	createdAccount, err := qtx.CreateAccount(ctx, sqlc.CreateAccountParams{
		Email:             email.String(),
		PasswordHash:      passwordHash.String(),
		PasswordAlgorithm: algorithm.String(),
	})
	if err != nil {
		r.logger.Error("Failed to create account", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}

	// delete temporary account
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(tmpAccount.Id),
		Valid: true,
	}
	_, err = qtx.DeleteTemporaryAccount(ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to delete temporary account", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}

	// --- commit the transaction
	err = transaction.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, err
	}

	// generate JWT token for the new account
	account_id := domain.AccountID(createdAccount.ID.Bytes)
	jwtToken, err := jwt.GenerateJwt(account_id)
	if err != nil {
		r.logger.Error("Failed to generate JWT", slog.String("account_id", account_id.String()), slog.String("error", err.Error()))
		return nil, err
	}

	jwtTokenStr := domain.JwtToken(*jwtToken)
	return &jwtTokenStr, nil
}

func (r *AccountRepository) Delete(ctx context.Context, id domain.AccountID) error {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}

	_, err := r.sql.DeleteAccount(ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to delete account", slog.String("account_id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}
