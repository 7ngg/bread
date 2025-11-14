package pages

import (
	"fmt"
	"net/http"
)

type NavbarItem struct {
	Name string
	Url  string
}

type IndexProps struct {
	NavbarItems []NavbarItem
}

func (handler *PagesHandler) RenderIndex(w http.ResponseWriter, r *http.Request) {
	// TODO: handle errors somehow
	products, err := handler.productService.GetAllProducts(r.Context(), 1, 10)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching products: %v", err), http.StatusInternalServerError)
		return
	}
	if getSessionID(r) == "" {
		setSessionCookie(w)
	}
	err = handler.Render(w, "index.html", products)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

