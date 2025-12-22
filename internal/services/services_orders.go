package services

import (
	"context"
	"fmt"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/pkg/validator"
	"github.com/akarashov/gofermart/internal/storage"
)

type orderService struct {
	orderRepo storage.OrderRepository
}

func NewOrderService(orderRepo storage.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

// Upload new order for processing
// Validates order number using Luhn algorithm
// Checks if order already exists for this user or another user
// Creates new order with status NEW if valid and not existing
func (s *orderService) UploadOrder(ctx context.Context, userID, orderNumber string) error {
	if !validator.ValidateLuhn(orderNumber) {
		return ErrInvalidOrderNumber
	}
	existingOrder, _ := s.orderRepo.GetOrderByNumber(ctx, userID, orderNumber)

	if existingOrder.Number != "" {
		return ErrOrderAlreadyUploadedByUser
	}
	existingByOther, err := s.orderRepo.GetOrderByNumberAnyUser(ctx, orderNumber)
	if err == nil && existingByOther.Number != "" && existingByOther.UserID != userID {
		return ErrOrderAlreadyUploadedByAnotherUser
	}
	order := models.Order{
		UserID: userID,
		Number: orderNumber,
		Status: models.OrderStatusNew,
	}
	err = s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

// GetUserOrders - get all orders for specific user
// Returns ErrNoOrdersForUser if no orders found
func (s *orderService) GetUserOrders(ctx context.Context, userID string) ([]models.OrderResponse, error) {
	orders, err := s.orderRepo.GetOrdersByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	if len(orders) == 0 {
		return nil, ErrNoOrdersForUser
	}
	return orders, nil
}
