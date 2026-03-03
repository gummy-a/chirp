package repository

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/domain/value_object"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type MediaRepository struct {
	logger slog.Logger
	sql    sqlc.Queries
	ctx    context.Context
}

func NewMediaRepository(logger slog.Logger, sql sqlc.Queries, ctx context.Context) MediaRepository {
	return MediaRepository{
		logger: logger,
		sql:    sql,
		ctx:    ctx,
	}
}

func (m *MediaRepository) SaveMetaDataToDB(metadata entity.MetaData, job entity.EncodeJob) error {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(job.OwnerAccountId),
		Valid: true,
	}

	meta, err := json.Marshal(metadata)
	if err != nil {
		m.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
		return err
	}

	_, err = m.sql.InsertMedia(m.ctx, sqlc.InsertMediaParams{
		OwnerAccountID:     pgtypeUUID,
		MimeType:           string(job.FileInfo.MimeType),
		OriginalFileName:   string(job.FileInfo.OriginalFileName),
		UnprocessedFileUrl: string(job.FileInfo.FileUrl),
		Metadata:           meta,
	})
	if err != nil {
		m.logger.Error("InsertMedia failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (m *MediaRepository) SaveFileToStorage(url value_object.FileUrl) error {
	// TODO: implement this
	return nil
}
