package pgxtransactor

import (
	"context"
)

type DB interface {
	Transactor
	QueryEngine
}

type TxRunner interface {
	RunRepeatableRead(ctx context.Context, fx func(dbCtx context.Context) error) error
	RunReadCommitted(ctx context.Context, fx func(dbCtx context.Context) error) error
}
