package repository

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/gummy_a/chirp/auth/internal/domain"
	"github.com/gummy_a/chirp/auth/internal/domain/entity"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type TemporaryAccountRepository struct {
	sql    *sqlc.Queries
	logger *slog.Logger
}

func NewTemporaryAccountRepository(q *sqlc.Queries, logger *slog.Logger) *TemporaryAccountRepository {
	return &TemporaryAccountRepository{sql: q, logger: logger}
}

func (r *TemporaryAccountRepository) Create(ctx context.Context, email domain.Email, passwordHash domain.PasswordHash, expiresAt domain.Timestamp) (*domain.NumberCode, *domain.TemporaryAccountID, error) {
	ts := pgtype.Timestamptz{
		Time:             time.Time(expiresAt),
		Valid:            true,
		InfinityModifier: pgtype.Finite,
	}

	numberCode := rand.Int31n(999999)
	token := domain.NumberCode(numberCode)

	tmpAccount, err := r.sql.CreateTemporaryAccount(ctx, sqlc.CreateTemporaryAccountParams{
		Email:        email.String(),
		PasswordHash: passwordHash.String(),
		ExpiresAt:    ts,
		NumberCode:   int32(numberCode),
	})
	if err != nil {
		r.logger.Error("Failed to create temporary account", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, nil, err
	}

	accountID := domain.TemporaryAccountID(tmpAccount.ID.Bytes)
	return &token, &accountID, nil
}

func (r *TemporaryAccountRepository) Delete(ctx context.Context, id *domain.TemporaryAccountID) error {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(*id),
		Valid: true,
	}

	_, err := r.sql.DeleteTemporaryAccount(ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to delete temporary account", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *TemporaryAccountRepository) FindById(ctx context.Context, id *domain.TemporaryAccountID) (*entity.TemporaryAccount, error) {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(*id),
		Valid: true,
	}

	dbTempAccount, err := r.sql.FindTemporaryAccountById(ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to find temporary account by id", slog.String("id", id.String()), slog.String("error", err.Error()))
		return nil, err
	}

	tempAccount := &entity.TemporaryAccount{
		Email:      domain.Email(dbTempAccount.Email),
		ExpiresAt:  domain.Timestamp(dbTempAccount.ExpiresAt.Time),
		Password:   domain.PasswordHash(dbTempAccount.PasswordHash),
		Id:         domain.TemporaryAccountID(dbTempAccount.ID.Bytes),
		NumberCode: domain.NumberCode(dbTempAccount.NumberCode),
	}
	return tempAccount, nil
}
