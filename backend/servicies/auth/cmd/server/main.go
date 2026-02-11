package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"log/slog"

	"github.com/gummy_a/chirp/auth/internal/adapter/handler"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/db"
	"github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/db/sqlc"
	repository "github.com/gummy_a/chirp/auth/internal/infrastructure/persistence/repository/impl"
	usecase "github.com/gummy_a/chirp/auth/internal/usecase/register_account"
)

func main() {
	ctx := context.Background()

	// 構造化ロギングのセットアップ
	opts := &slog.HandlerOptions{AddSource: true}
	jsonhandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(jsonhandler)

	// Infrastructure層：DB接続を確立
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	// Infrastructure層：Queriesオブジェクトを作成
	queries := sqlc.New(pool)

	// Repository層：Repository実装を生成
	accountRepo := repository.NewAccountRepository(pool, queries, logger)
	temporaryAccountRepo := repository.NewTemporaryAccountRepository(queries, logger)
	registrationSenderRepo := repository.NewRegistrationSenderRepository(logger)

	// UseCase層：ユースケースを生成
	SignupAccountUseCase := usecase.NewSignupAccountUseCase(accountRepo, temporaryAccountRepo)
	SignupTemporaryAccountUseCase := usecase.NewSignupTemporaryAccountUseCase(temporaryAccountRepo, registrationSenderRepo)

	// Adapter層：UseCaseと接合
	router := handler.NewSignupRouter(SignupTemporaryAccountUseCase, SignupAccountUseCase, logger)

	// APIサーバー待ち受け開始
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
