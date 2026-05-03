package external

import (
	"chirp/backend/router"
	authController "chirp/backend/services/auth/v1/internal/adapter/controller"
	handler "chirp/backend/services/auth/v1/internal/adapter/http"
	"chirp/backend/services/auth/v1/internal/infrastructure/email"
	"chirp/backend/services/auth/v1/internal/infrastructure/persistence/db"
	"chirp/backend/services/auth/v1/internal/infrastructure/persistence/db/sqlc"
	"chirp/backend/services/auth/v1/internal/infrastructure/persistence/repository"
	loginUsecase "chirp/backend/services/auth/v1/internal/usecase/login"
	signupUsecase "chirp/backend/services/auth/v1/internal/usecase/signup"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func checkEnvironmentVariables() {
	env := os.Getenv("AUTH_SERVICE_APP_ENV")
	if env == "" {
		log.Fatal("AUTH_SERVICE_APP_ENV environment variable is not set")
	}

	jwtSecretKey := os.Getenv("AUTH_SERVICE_JWT_SECRET_KEY")
	if env == "production" && jwtSecretKey == "" {
		log.Fatal("AUTH_SERVICE_JWT_SECRET_KEY is not set")
	}

	url := os.Getenv("AUTH_SERVICE_DATABASE_URL")
	if env == "production" && url == "" {
		log.Fatal("AUTH_SERVICE_DATABASE_URL is not set")
	}

	region := os.Getenv("AUTH_SERVICE_AWS_SES_REGION")
	if env == "production" && region == "" {
		log.Fatal("AUTH_SERVICE_AWS_SES_REGION is not set")
	}

	from := os.Getenv("AUTH_SERVICE_AWS_SES_FROM_ADDRESS")
	if env == "production" && from == "" {
		log.Fatal("AUTH_SERVICE_AWS_SES_FROM_ADDRESS is not set")
	}

	front := os.Getenv("AUTH_SERVICE_FRONTEND_URL")
	if env == "production" && front == "" {
		log.Fatal("AUTH_SERVICE_FRONTEND_URL is not set")
	}
}

type CreateAuthServiceResult struct {
	Controller []router.Handler
	Pool       *pgxpool.Pool
}

func CreateAuthService() (*CreateAuthServiceResult, error) {
	checkEnvironmentVariables()
	ctx := context.Background()

	// setup logger
	opts := &slog.HandlerOptions{AddSource: true}
	jsoncontroller := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsoncontroller)

	// Infrastructure layer: create DB pool
	pool, err := db.NewPool(ctx)
	if err != nil {
		return nil, err
	}

	// Infrastructure layer: create SQLC sql
	sql := sqlc.New(pool)

	// Infrastructure layer: create email sender
	cfg, err := email.NewAwsConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to create AWS config: %v", err)
	}
	sender := email.NewEmailSender(*cfg, ctx)

	// Repository layer: create repositories
	accountRepo := repository.NewAccountRepository(pool, sql, logger, &ctx)
	temporaryAccountRepo := repository.NewTemporaryAccountRepository(sql, logger, &ctx)

	// UseCase layer: create use cases
	signupAccountUseCase := signupUsecase.NewSignupAccountUseCase(accountRepo, temporaryAccountRepo)
	signupTemporaryAccountUseCase := signupUsecase.NewSignupTemporaryAccountUseCase(temporaryAccountRepo, sender)
	loginUseCase := loginUsecase.NewLoginAccountUseCase(accountRepo)

	// Adapter layer: create controller
	controller := handler.NewAuthHandler(&authController.AuthControllerUseCases{
		Signup:    signupAccountUseCase,
		Login:     loginUseCase,
		TmpSignup: signupTemporaryAccountUseCase,
	}, logger)

	return &CreateAuthServiceResult{
		Controller: controller,
		Pool:       pool,
	}, nil
}
