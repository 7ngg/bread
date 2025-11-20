package services

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type BasketService struct {
	client *redis.Client
}

func NewBasketService(client *redis.Client) *BasketService {
	return &BasketService{
		client: client,
	}
}

type Basket struct {
	Items map[int]int `json:"items" redis:"items"`
}

func (bs *BasketService) UpdateBasketCount(ctx context.Context, sessionID string, itemID, delta int) (id int, count int, err error) {
	basket, err := bs.GetBasket(ctx, sessionID)	
	if err == redis.Nil {
		basket = Basket{ Items: make(map[int]int) }
	} else if err != nil {
		return -1, -1, err
	}

	if basket.Items[itemID] <= 0 && delta <= 0 {
		return itemID, basket.Items[itemID], err
	}

	basket.Items[itemID] += delta

	if basket.Items[itemID] == 0 {
		delete(basket.Items, itemID)
	}

	if err = bs.client.JSONSet(ctx, sessionID, "$", basket).Err();  err != nil {
		return -1, -1, err
	}

	return itemID, basket.Items[itemID], nil
}

func (bs *BasketService) AddItemToBasket(ctx context.Context, sessionID string, itemID int) (id int, count int, err error) {
	return bs.UpdateBasketCount(ctx, sessionID, itemID, 1)
}

func (bs *BasketService) RemoveItemFromBasket(ctx context.Context, sessionID string, itemID int) (id int, count int, err error) {
	return bs.UpdateBasketCount(ctx, sessionID, itemID, -1)
}

func (bs *BasketService) GetBasket(ctx context.Context, sessionID string) (Basket, error) {
	basket, err := bs.client.JSONGet(ctx, sessionID).Result()	
	if err != nil {
		return Basket{}, err
	}

	var obj Basket
	err = json.Unmarshal([]byte(basket), &obj)
	return obj, err
}
