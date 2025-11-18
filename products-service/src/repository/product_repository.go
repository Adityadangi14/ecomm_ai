package repository

import (
	"context"

	"github.com/Adityadangi14/ecomm_ai/pkg/WDB"
	"github.com/Adityadangi14/ecomm_ai/products-service/src/models"
)

type ProductRepository interface {
	SaveProduct(ctx context.Context, data models.ProductsModel) error
	NearSearchProducts(ctx context.Context, query string) (models.ProductsModel, error)
	DeleteAllProducts(ctx context.Context) error
}

type prodRepo struct {
	WDB *WDB.WDB
}

func NewProductRepository(wdb *WDB.WDB) ProductRepository {
	return &prodRepo{WDB: wdb}
}

func (p *prodRepo) SaveProduct(ctx context.Context, data models.ProductsModel) error {
	return nil
}

func (p *prodRepo) NearSearchProducts(ctx context.Context, query string) (models.ProductsModel, error) {
	return models.ProductsModel{}, nil
}

func (p *prodRepo) DeleteAllProducts(ctx context.Context) error {
	return nil
}
