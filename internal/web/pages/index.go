package pages

import (
	"fmt"
	"net/http"

	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/lib"
	"github.com/7ngg/bread/internal/services"
)

type IndexProduct struct {
	ID          int32
	Name        string
	Ingredients string
	Price       float64
	ImgUrl      string
	Count       int
}

func NewIndexProduct(product db.Product, quantity int) *IndexProduct {
	return &IndexProduct{
		ID:          product.ID,
		Name:        product.Name,
		Ingredients: product.Ingredients,
		Price:       product.Price,
		ImgUrl:      product.ImgUrl,
		Count:       quantity,
	}
}

func merge(entities lib.PaginatedList[db.Product], basket services.Basket) []IndexProduct {
	products := make([]IndexProduct, 0, len(entities.Items))
	for _, p := range entities.Items {
		products = append(products, *NewIndexProduct(p, basket.Items[int(p.ID)]))
	}

	return products
}

func (handler *PagesHandler) RenderIndex(w http.ResponseWriter, r *http.Request) {
	// TODO: handle errors somehow
	products, err := handler.productService.GetAllProducts(r.Context(), 1, 10)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching products: %v", err), http.StatusInternalServerError)
		return
	}

	sessionID := getSessionID(r)
	if sessionID == "" {
		setSessionCookie(w, lib.RandomString(16))
	}

	// TODO: err
	basket, err := handler.basketService.GetBasket(r.Context(), sessionID)
	err = handler.Render(w, "index.html", merge(products, basket))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
