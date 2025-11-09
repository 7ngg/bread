package web

import (
	"net/http"
	"strconv"

	"github.com/7ngg/bread/internal/lib"
)

func (handler *WebApp) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	page, _ := strconv.Atoi(params.Get("page"))
	pageSize, _ := strconv.Atoi(params.Get("page-size"))

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 20
	}

	products, err := handler.productService.GetAllProducts(r.Context(), int32(page), int32(pageSize))
	if err != nil {
		lib.ResponseWithJson(w, http.StatusInternalServerError, "unable to retrieve products")
		return
	}

	lib.ResponseWithJson(w, http.StatusOK, products)
}

