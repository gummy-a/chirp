package repository

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"chirp/backend/services/auth/v1/internal/domain/entity"
	"chirp/backend/services/auth/v1/internal/domain/value_object"
	"chirp/backend/services/auth/v1/internal/infrastructure/persistence/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type TemporaryAccountRepository struct {
	sql    *sqlc.Queries
	logger *slog.Logger
	ctx    *context.Context
}

func NewTemporaryAccountRepository(q *sqlc.Queries, logger *slog.Logger, ctx *context.Context) *TemporaryAccountRepository {
	return &TemporaryAccountRepository{sql: q, logger: logger, ctx: ctx}
}

func (r *TemporaryAccountRepository) Create(email value_object.Email, passwordHash value_object.PasswordHash, expiresAt value_object.Timestamp) (*value_object.NumberCode, *value_object.TemporaryAccountID, error) {
	ts := pgtype.Timestamptz{
		Time:             time.Time(expiresAt),
		Valid:            true,
		InfinityModifier: pgtype.Finite,
	}

	numberCode := rand.IntN(900000) + 100000
	token := value_object.NumberCode(numberCode)

	tmpAccount, err := r.sql.CreateTemporaryAccount(*r.ctx, sqlc.CreateTemporaryAccountParams{
		Email:        email.String(),
		PasswordHash: passwordHash.String(),
		ExpiresAt:    ts,
		NumberCode:   int32(numberCode),
	})
	if err != nil {
		r.logger.Error("Failed to create temporary account", slog.String("email", email.String()), slog.String("error", err.Error()))
		return nil, nil, err
	}

	accountID := value_object.TemporaryAccountID(tmpAccount.ID.Bytes)
	return &token, &accountID, nil
}

func (r *TemporaryAccountRepository) Delete(id value_object.TemporaryAccountID) error {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}

	_, err := r.sql.DeleteTemporaryAccount(*r.ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to delete temporary account", slog.String("id", id.String()), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *TemporaryAccountRepository) FindById(id value_object.TemporaryAccountID) (*entity.TemporaryAccount, error) {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(id),
		Valid: true,
	}

	dbTempAccount, err := r.sql.FindTemporaryAccountById(*r.ctx, pgtypeUUID)
	if err != nil {
		r.logger.Error("Failed to find temporary account by id", slog.String("id", id.String()), slog.String("error", err.Error()))
		return nil, err
	}

	tempAccount := &entity.TemporaryAccount{
		Email:      value_object.Email(dbTempAccount.Email),
		ExpiresAt:  value_object.Timestamp(dbTempAccount.ExpiresAt.Time),
		Password:   value_object.PasswordHash(dbTempAccount.PasswordHash),
		Id:         value_object.TemporaryAccountID(dbTempAccount.ID.Bytes),
		NumberCode: value_object.NumberCode(dbTempAccount.NumberCode),
	}
	return tempAccount, nil
}
