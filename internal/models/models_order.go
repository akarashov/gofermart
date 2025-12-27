package models

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
)

var (
	ErrEmptyOrderNumber     = errors.New("empty order number")
	ErrInvalidOrderNumber   = errors.New("invalid order number")
	ErrInvalidWithdrawalSum = errors.New("invalid withdrawal sum")
)

// Order модель заказа
type Order struct {
	ID         string          `json:"-" db:"id"`
	UserID     string          `json:"-" db:"user_id"`
	Number     string          `json:"number" db:"order_number"`
	Status     OrderStatus     `json:"order_status" db:"order_status"`
	Accrual    decimal.Decimal `json:"accrual,omitempty" db:"accrual,omitempty"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UploadedAt time.Time       `json:"uploaded_at,omitempty" db:"uploaded_at,omitempty"`
}

type OrderUploadRequest struct {
	Number string `json:"order_number"`
}

// OrderResponse ответ при запросе списка заказов
type OrderResponse struct {
	Number     string      `json:"order_number"`
	Status     OrderStatus `json:"order_status"`
	Accrual    float32     `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

// OrderUpdate обновление заказа
type OrderUpdate struct {
	Number  string          `db:"order_number"`
	Status  string          `db:"order_status"`
	Accrual decimal.Decimal `db:"accrual,omitempty"`
}

type OrdersForProcessing struct {
	UserID string `json:"-" db:"user_id"`
	Number string `json:"number" db:"order_number"`
}
