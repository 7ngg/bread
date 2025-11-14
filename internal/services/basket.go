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

func (bs *BasketService) AddItemToBasket(ctx context.Context, sessionID string, itemID int) (int, error) {
	basket, err := bs.GetBasket(ctx, sessionID)	
	if err == redis.Nil {
		basket = Basket{ Items: make(map[int]int) }
	} else if err != nil {
		return -1, err
	}

	basket.Items[itemID]++

	if err = bs.client.JSONSet(ctx, sessionID, "$", basket).Err();  err != nil {
		return -1, err
	}

	return basket.Items[itemID], nil
}

func (bs *BasketService) RemoveItemToBasket(ctx context.Context, sessionID string, itemID int) (int, error) {
	basket, err := bs.GetBasket(ctx, sessionID)	
	if err == redis.Nil {
		basket = Basket{ Items: make(map[int]int) }
	} else if err != nil {
		return -1, err
	}

	basket.Items[itemID]--

	if err = bs.client.JSONSet(ctx, sessionID, "$", basket).Err();  err != nil {
		return -1, err
	}

	return basket.Items[itemID], nil
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
