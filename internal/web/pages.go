package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	if getSessionID(r) == "" {
		setSessionCookie(w)
	}
	err = t.Render(w, "index.html", products)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func setSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: lib.RandomString(16),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func getSessionID(r *http.Request) string {
	c, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	return c.Value
}

func (t *WebApp) RenderAddItemToBasket(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		lib.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	sessionID := getSessionID(r)
	if sessionID == "" {
		lib.RespondWithError(w, http.StatusBadRequest, "No session ID found")
		return
	}

	err = t.basketService.AddItemToBasket(r.Context(), sessionID, productID)
	if err != nil {
		lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t.Render(w, "basket", nil)
}
