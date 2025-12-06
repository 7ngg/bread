package services

import (
	"context"
	"fmt"

	"github.com/7ngg/bread/internal/db"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, arg db.CreateOrderParams) (db.Order, error)
}

type IProductRepository interface {
	GetProductById(ctx context.Context, id int32) (db.Product, error)
}

type IOrderItemRepository interface {
	InsertOrderItems(ctx context.Context, arg []db.InsertOrderItemsParams) (int64, error)
}

type OrderService struct {
	orderRepository     IOrderRepository
	productRepository   IProductRepository
	orderItemRepository IOrderItemRepository
	userService *UserService
}

func NewOrderService(
	orderRepository IOrderRepository,
	productRepository IProductRepository,
	orderItemRepository IOrderItemRepository,
	userService *UserService,
) *OrderService {
	return &OrderService{
		orderRepository:     orderRepository,
		productRepository:   productRepository,
		orderItemRepository: orderItemRepository,
		userService: userService,
	}
}

type OrderItem struct {
	ID       int32
	Quantity int32
}

func (o *OrderService) NewOrder(ctx context.Context, phone, name string, items []OrderItem) (db.Order, error) {
	totalPrice, err := o.calculateTotalPrice(ctx, items)
	if err != nil {
		return db.Order{}, err
	}
	fmt.Printf("total price: %f", totalPrice)

	user, err := o.userService.EnsureUserExists(ctx, phone, name)
	if err != nil {
		return db.Order{}, err
	}

	order, err := o.orderRepository.CreateOrder(ctx, db.CreateOrderParams{
		UserID:     user.ID,
		TotalPrice: totalPrice,
	})

	if err != nil {
		return db.Order{}, err
	}

	params := make([]db.InsertOrderItemsParams, 0, len(items))
	for _, item := range items {
		params = append(params, db.InsertOrderItemsParams{
			OrderID:   order.ID,
			ProductID: item.ID,
			Quantity:  int32(item.Quantity),
		})
	}

	if rowsAffected, err := o.orderItemRepository.InsertOrderItems(ctx, params); err != nil || rowsAffected == 0 {
		return db.Order{}, nil
	}

	return order, nil
}

func (o *OrderService) calculateTotalPrice(ctx context.Context, items []OrderItem) (float64, error) {
	price := 0.0

	for _, item := range items {
		prod, err := o.productRepository.GetProductById(ctx, item.ID)
		if err != nil {
			fmt.Printf("could not find product %d", item.ID)
			return 0.0, err
		}

		price +=  prod.Price * float64(item.Quantity)
	}

	return price, nil
}
