package services

import (
	"context"
	"math"

	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/lib"
)

type ProductService struct {
	ProductsGetter IProductsGetter
}

type IProductsGetter interface {
	GetProducts(context context.Context, params db.GetProductsParams) ([]db.Product, error)
	ProductsCount(context context.Context) (int64, error)
	GetProductById(context context.Context, id int32) (db.Product, error)
}

func NewProductService(pg IProductsGetter) *ProductService {
	return &ProductService{
		ProductsGetter: pg,
	}
}

func (ps *ProductService) GetAllProducts(context context.Context, page, pageSize int32) (lib.PaginatedList[db.Product], error) {
	params := db.GetProductsParams{Limit: pageSize, Offset: (page - 1) * pageSize}
	products, err := ps.ProductsGetter.GetProducts(context, params)
	if err != nil {
		return lib.PaginatedList[db.Product]{}, err
	}

	totalItems, err := ps.ProductsGetter.ProductsCount(context)
	if err != nil {
		return lib.PaginatedList[db.Product]{}, err
	}

	response := lib.PaginatedList[db.Product]{
		Items:      products,
		TotalPages: int32(math.Ceil(float64(totalItems) / float64(pageSize))),
		Page:       page,
	}

	return response, nil
}
