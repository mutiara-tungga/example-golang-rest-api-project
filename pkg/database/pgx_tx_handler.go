package database

import (
	"context"
	"golang-rest-api/internal/model"
	pkgErr "golang-rest-api/pkg/error"
	"golang-rest-api/pkg/log"

	"github.com/jackc/pgx/v5"
)

type PgxTxFn func(context.Context, pgx.Tx) error

type TxHandler interface {
	WithTransaction(context.Context, PgxTxFn) error
}

type pgxTxHandler struct {
	db IPostgres
}

func NewTxHandler(db IPostgres) TxHandler {
	return &pgxTxHandler{db: db}
}

func (th *pgxTxHandler) WithTransaction(ctx context.Context, fn PgxTxFn) (err error) {
	tx, err := th.db.Begin(ctx)
	if err != nil {
		log.Error(ctx, "failed to begin transaction ", err)
		return pkgErr.NewCustomErrWithOriginalErr(model.ErrorExecQuery, err)
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Error(ctx, "failed to rollback transaction ", rollbackErr)
			}
			panic(p)
		}

		if err != nil {
			// something went wrong, rollback
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Error(ctx, "failed to rollback transaction ", rollbackErr)
			}
			return
		}

		// all good, commit
		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			log.Error(ctx, "failed to commit transaction", err)
			err = pkgErr.NewCustomErrWithOriginalErr(model.ErrorExecQuery, commitErr)
		}
	}()

	err = fn(ctx, tx)
	return
}
