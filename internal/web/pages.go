package web

import (
	"fmt"
	"net/http"

	"github.com/7ngg/bread/internal/lib"
)

type NavbarItem struct {
	Name string
	Url  string
}

type IndexProps struct {
	NavbarItems []NavbarItem
}

func (t *WebApp) RenderIndex(w http.ResponseWriter, r *http.Request) {
	// TODO: handle errors somehow
	products, err := t.productService.GetAllProducts(r.Context(), 1, 10)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching products: %v", err), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: lib.RandomString(16),
	})
	err = t.Render(w, "index.html", products)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
