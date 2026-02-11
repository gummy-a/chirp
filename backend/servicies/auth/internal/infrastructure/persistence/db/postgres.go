package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
}
