package database

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
)

type IPostgres interface {
	Begin(ctx context.Context) (ITx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (ITx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type ITx interface {
	Begin(ctx context.Context) (ITx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
