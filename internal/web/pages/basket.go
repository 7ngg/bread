package pages

import (
	"strconv"
	"net/http"
	
	"github.com/7ngg/bread/internal/lib"
)

type CounterProps struct {
	Count int
}

func (handler *PagesHandler) RenderPlus(w http.ResponseWriter, r *http.Request) {
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

	_, err = handler.basketService.AddItemToBasket(r.Context(), sessionID, productID)
	if err != nil {
		lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Render(w, "counter", nil)
}
