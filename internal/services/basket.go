package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

func (bs *BasketService) AddItemToBasket(context context.Context, sessionID string, itemID int) error {
	var basket Basket
	
	data, err := bs.client.Get(context, sessionID).Bytes()
	if err == nil {
		if err := json.Unmarshal(data, &basket); err != nil {
			return fmt.Errorf("failed to unmarshal basket: %w", err)
		}
	} else if err != redis.Nil {
		return fmt.Errorf("failed to get basket: %w", err)
	}
	
	if basket.Items == nil {
		basket.Items = make(map[int]int)
	}
	
	basket.Items[itemID]++
	
	jsonData, err := json.Marshal(basket)
	if err != nil {
		return fmt.Errorf("failed to marshal basket: %w", err)
	}
	
	return bs.client.Set(context, sessionID, jsonData, 24*time.Hour).Err()
}
