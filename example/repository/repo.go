package repository

import (
	"context"

	"github.com/VladimirStepanov/pgxtransactor"
	"github.com/VladimirStepanov/pgxtransactor/example/model"
	"github.com/georgysavva/scany/pgxscan"
)

type Repository interface {
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
}

type repository struct {
	qp pgxtransactor.QueryEngineProvider
}

func New(qp pgxtransactor.QueryEngineProvider) *repository {
	return &repository{
		qp: qp,
	}
}

func (r *repository) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	db := r.qp.GetQueryEngine(ctx)

	query := `SELECT * FROM products WHERE id=$1 FOR UPDATE`

	var product model.Product

	rows, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(&product, rows); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) UpdateProduct(ctx context.Context, product *model.Product) error {
	db := r.qp.GetQueryEngine(ctx)
	query := `UPDATE products SET name=$1, price=$2 WHERE id=$3`
	_, err := db.Exec(ctx, query, product.Name, product.Price, product.ID)
	return err
}
