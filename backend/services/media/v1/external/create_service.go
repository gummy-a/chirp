package external

import (
	"chirp/backend/router"
	"chirp/backend/services/media/v1/internal/adapter/controller"
	"chirp/backend/services/media/v1/internal/adapter/http"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/db"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/db/sqlc"
	"chirp/backend/services/media/v1/internal/infrastructure/persistence/repository"
	"chirp/backend/services/media/v1/internal/infrastructure/redis"
	"chirp/backend/services/media/v1/internal/usecase"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func checkEnvironmentVariables() {
	env := os.Getenv("MEDIA_SERVICE_APP_ENV")
	if env == "" {
		log.Fatal("MEDIA_SERVICE_APP_ENV environment variable is not set")
	}

	jwtSecretKey := os.Getenv("MEDIA_SERVICE_JWT_SECRET_KEY")
	if env == "production" && jwtSecretKey == "" {
		log.Fatal("MEDIA_SERVICE_JWT_SECRET_KEY is not set")
	}

	url := os.Getenv("MEDIA_SERVICE_DATABASE_URL")
	if env == "production" && url == "" {
		log.Fatal("MEDIA_SERVICE_DATABASE_URL is not set")
	}

	redis := os.Getenv("MEDIA_SERVICE_REDIS_URL")
	if env == "production" && redis == "" {
		log.Fatal("MEDIA_SERVICE_REDIS_URL is not set")
	}

	s3 := os.Getenv("MEDIA_SERVICE_AWS_S3_BUCKET_NAME")
	if env == "production" && s3 == "" {
		log.Fatal("MEDIA_SERVICE_AWS_S3_BUCKET_NAME is not set")
	}

	region := os.Getenv("MEDIA_SERVICE_AWS_S3_REAGION")
	if env == "production" && region == "" {
		log.Fatal("MEDIA_SERVICE_AWS_S3_REAGION is not set")
	}
}

type CreateMediaServiceResult struct {
	Controller []router.Handler
	Pool       *pgxpool.Pool
}

func CreateMediaService() (*CreateMediaServiceResult, error) {
	checkEnvironmentVariables()
	ctx := context.Background()

	// setup logger
	opts := &slog.HandlerOptions{AddSource: true}
	jsoncontroller := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsoncontroller)

	// Infrastructure layer: create DB dbPool
	dbPool, err := db.NewConnectionPool(ctx)
	if err != nil {
		return nil, err
	}

	// Infrastructure layer: create database object
	sql := sqlc.New(dbPool)

	// Infrastructure layer: create queue handler
	queue := redis.NewQueueHandler(ctx, logger)

	// Repository layer: create repositories
	s3, err := repository.NewS3Client(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}
	mediaRepository := repository.NewMediaRepository(logger, sql, s3, ctx)

	// UseCase layer: create use cases
	mediaControlUseCase := usecase.NewMediaUploadUseCase(mediaRepository, queue)
	encodeStatusUsecase := usecase.NewMediaEncodeStatusUseCase(queue)

	// Adapter layer: create HTTP controllers and router
	mediaHandler := handler.NewMediaHandler(&mediaController.MediaControllerUseCases{
		Upload: *mediaControlUseCase,
		Status: *encodeStatusUsecase,
	}, logger)

	return &CreateMediaServiceResult{
		Controller: mediaHandler,
		Pool:       dbPool,
	}, nil
}
