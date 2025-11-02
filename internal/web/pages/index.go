package pages

import (
	"net/http"
)

type NavbarItem struct {
	Name string
	Url  string
}

type IndexProps struct {
	NavbarItems []NavbarItem
}

func (t *Templates) RenderIndex(w http.ResponseWriter, r *http.Request) {
	err := t.Render(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
