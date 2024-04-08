package transaction

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...any) pgx.Row
}

type TxFunc func(context.Context, func(tx pgx.Tx) error) error

func Tx(p *pgxpool.Pool, log *slog.Logger) TxFunc {
	return func(ctx context.Context, cb func(pgx.Tx) error) error {
		tx, err := p.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if err := cb(tx); err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Error("failed to rollback transaction", slog.String("err", rbErr.Error()))
			}

			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}
