package pages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/7ngg/bread/internal/lib"
)

type CounterProps struct {
	ItemID int
	Count  int
}

var (
	increment = "inc"
	decrement = "dec"
)

func (handler *PagesHandler) RenderPlus(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		lib.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	action := r.FormValue("action")

	sessionID := getSessionID(r)
	if sessionID == "" {
		lib.RespondWithError(w, http.StatusBadRequest, "No session ID found")
		return
	}

	var itemID int
	var updatedCount int
	switch action {
	case increment:
		itemID, updatedCount, err = handler.basketService.AddItemToBasket(r.Context(), sessionID, productID)
	case decrement:
		itemID, updatedCount, err = handler.basketService.RemoveItemFromBasket(r.Context(), sessionID, productID)
	default:
		err = errors.New("invalid action")
	}
	if err != nil {
		lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.Render(w, "counter", CounterProps{itemID, updatedCount})
}
