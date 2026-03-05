package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gummy_a/chirp/media/internal/domain/entity"
	"github.com/gummy_a/chirp/media/internal/infrastructure/persistence/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type MediaRepository struct {
	logger slog.Logger
	sql    sqlc.Queries
	ctx    context.Context
	s3     s3.Client
}

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	region := os.Getenv("MEDIA_SERVICE_S3_REAGION")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}

func NewMediaRepository(logger slog.Logger, sql sqlc.Queries, s3 s3.Client, ctx context.Context) MediaRepository {
	return MediaRepository{
		logger: logger,
		sql:    sql,
		ctx:    ctx,
		s3:     s3,
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
		MimeType:           string(job.UploadedFileInfo.MimeType),
		OriginalFileName:   string(job.UploadedFileInfo.OriginalFileName),
		UnprocessedFileUrl: string(job.UploadedFileInfo.FileUrl),
		Metadata:           meta,
	})
	if err != nil {
		m.logger.Error("InsertMedia failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (m *MediaRepository) SaveFileToStorage(file entity.UploadedFileInfo) error {
	bucketName := os.Getenv("MEDIA_SERVICE_S3_BUCKET_NAME")
	body, err := os.Open(string(file.OriginalFileName))
	if err != nil {
		m.logger.Error("os.Open failed", slog.String("error", err.Error()))
		return err
	}
	defer body.Close()

	_, err = m.s3.PutObject(m.ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(string(file.FileUrl)),
		Body:   body,
	})

	if err != nil {
		m.logger.Error("s3.PutObject failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}
