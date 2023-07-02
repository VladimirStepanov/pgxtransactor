package pgxtransactor

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type txKey string

const key = txKey("tx")

type transactionManager struct {
	db DB
}

func NewTransactionManager(db DB) *transactionManager {
	return &transactionManager{db: db}
}

func (tm *transactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.db
}

func rollback(ctx context.Context, tx pgx.Tx, err error) error {
	if rerr := tx.Rollback(ctx); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}

func (tm *transactionManager) runTx(ctx context.Context, txLevel pgx.TxIsoLevel, fx func(ctxTX context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: txLevel,
		})
	if err != nil {
		return err
	}
	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		return rollback(ctx, tx, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return rollback(ctx, tx, err)
	}
	return nil
}

func (tm *transactionManager) RunReadCommitted(ctx context.Context, fx func(dbCtx context.Context) error) error {

	return tm.runTx(ctx, pgx.ReadCommitted, fx)
}

func (tm *transactionManager) RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error {
	return tm.runTx(ctx, pgx.RepeatableRead, fx)
}
