package pages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/7ngg/bread/internal/lib"
	"github.com/7ngg/bread/internal/services"
	"github.com/redis/go-redis/v9"
)

type CounterProps struct {
	ID int
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

type CheckoutItem struct {
	Name       string
	Quantity   int
	Price      float64
	TotalPrice float64
}

type CheckoutProps struct {
	TotalPrice float64
	Items      []CheckoutItem
}

func (handler *PagesHandler) RenderCheckout(w http.ResponseWriter, r *http.Request) {
	sessionID := getSessionID(r)
	if sessionID == "" {
		lib.RespondWithError(w, http.StatusBadRequest, "No session ID found")
		return
	}
	basket, err := handler.basketService.GetBasket(r.Context(), sessionID)
	if err == redis.Nil {
		basket = services.Basket{}
	} else if err != nil {
		lib.RespondWithError(w, http.StatusInternalServerError, "unable to retrieve basket")
		return
	}

	props := CheckoutProps{TotalPrice: 0.0, Items: make([]CheckoutItem, 0, len(basket.Items))}
	for id, count := range basket.Items {
		prod, err := handler.productService.ProductsGetter.GetProductById(r.Context(), int32(id))
		if err != nil {
			lib.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		price, _ := strconv.ParseFloat(prod.Price, 8)
		props.Items = append(props.Items, CheckoutItem{
			Name: prod.Name,
			Quantity: count,
			Price: price,
			TotalPrice: float64(count) * price,
		})
		props.TotalPrice += float64(count) * price
	}

	handler.Render(w, "checkout", props)
}
