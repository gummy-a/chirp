package repository

import (
	"context"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaRepository struct {
	db     *pgxpool.Pool
	sql    *sqlc.Queries
	logger *slog.Logger
}

func NewMediaRepository(db *pgxpool.Pool, q *sqlc.Queries, logger *slog.Logger) *MediaRepository {
	return &MediaRepository{db: db, sql: q, logger: logger}
}

func (r *MediaRepository) Save(ctx context.Context, files *[]entity.OriginalFileInfo, owner_account_id *domain.OwnerAccountId) error {
	// TODO: implement enqueue
	//TODO: imple Save()
	return nil
}
