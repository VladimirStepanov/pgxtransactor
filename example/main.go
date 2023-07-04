package main

import (
	"context"
	"log"

	"github.com/VladimirStepanov/pgxtransactor"
	"github.com/VladimirStepanov/pgxtransactor/example/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbAddr = "postgres://user:password@127.0.0.1:5433/example?sslmode=disable"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), dbAddr)
	if err != nil {
		log.Fatal("pxpool.Connect error", err)
	}
	defer pool.Close()

	txManager := pgxtransactor.NewTransactionManager(pool)

	ctx := context.Background()

	var repo repository.Repository = repository.New(txManager)

	err = txManager.RunRepeatableRead(ctx, func(dbCtx context.Context) error {
		product, err := repo.GetProduct(dbCtx, 1)
		if err != nil {
			return err
		}
		// some data processing
		product.Price += 10000
		return repo.UpdateProduct(dbCtx, product)
	})

	if err != nil {
		log.Fatal(err)
	}

}
