package database

import (
	"context"
	"net/http"

	pkgErr "golang-rest-api/pkg/error"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	RecordNotFound = pkgErr.NewCustomError("Record Not Found", "RECORD_NOT_FOUND", http.StatusNotFound)
)

type IPostgres interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row

	Get(ctx context.Context, destination any, query string, args ...any) error
	Select(ctx context.Context, destination any, query string, args ...any) error
}
