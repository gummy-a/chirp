package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"chirp/backend/services/media/v1/internal/domain/entity"
	"chirp/backend/services/media/v1/internal/domain/value_object"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/db/sqlc"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
)

type MediaRepository struct {
	logger *slog.Logger
	sql    *sqlc.Queries
	s3     *s3.Client
	ctx    context.Context
}

func NewS3Client(ctx context.Context) (*s3.Client, error) {
	region := os.Getenv("MEDIA_SERVICE_AWS_S3_REAGION")

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}

func NewMediaRepository(logger *slog.Logger, sql *sqlc.Queries, s3 *s3.Client, ctx context.Context) *MediaRepository {
	return &MediaRepository{
		logger: logger,
		sql:    sql,
		ctx:    ctx,
		s3:     s3,
	}
}

func (m *MediaRepository) Save(media *entity.Media) (*value_object.MediaId, error) {
	pgtypeUUID := pgtype.UUID{
		Bytes: [16]byte(media.MediaInfo.OwnerAccountId),
		Valid: true,
	}

	meta, err := json.Marshal(media.Metadata)
	if err != nil {
		m.logger.Error("json.Unmarshal failed", slog.String("error", err.Error()))
		return nil, err
	}

	ret, err := m.sql.InsertMedia(m.ctx, sqlc.InsertMediaParams{
		OwnerAccountID:   pgtypeUUID,
		MimeType:         string(media.MediaInfo.UploadedFile.MimeType),
		OriginalFileName: string(media.MediaInfo.UploadedFile.OriginalFileName),
		FileUrl:          string(media.MediaInfo.UploadedFile.FileUrl),
		Metadata:         meta,
	})
	if err != nil {
		m.logger.Error("InsertMedia failed", slog.String("error", err.Error()))
		return nil, err
	}

	mediaId := value_object.MediaId(ret.ID.Bytes)
	return &mediaId, nil
}

func (m *MediaRepository) SaveFileToStorage(path value_object.RealPath, url value_object.FileUrl) error {
	env := os.Getenv("MEDIA_SERVICE_APP_ENV")
	if env == "development" {
		// skip uploading in dev env
		return nil
	}

	bucketName := os.Getenv("MEDIA_SERVICE_AWS_S3_BUCKET_NAME")
	body, err := os.Open(string(path))
	if err != nil {
		m.logger.Error("os.Open failed", slog.String("error", err.Error()))
		return err
	}
	defer body.Close()

	_, err = m.s3.PutObject(m.ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(string(url)),
		Body:   body,
	})

	if err != nil {
		m.logger.Error("s3.PutObject failed", slog.String("error", err.Error()))
		return err
	}

	m.logger.Info("saved \n", slog.String("path", string(path)), slog.String("url", string(url)))
	return nil
}

func (m *MediaRepository) Delete(path value_object.RealPath) error {
	return os.Remove(path.String())
}
