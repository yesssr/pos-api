package configuration

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
  dsn := os.Getenv("DATABASE_URL")
  if dsn == "" {
    return nil, fmt.Errorf("DATABASE_URL not set")
  }

  pool, err := pgxpool.New(ctx, dsn)
  if err != nil {
    return nil, err
  }
  pool.BeginTx(ctx, pgx.TxOptions{});
  return pool, nil
}
