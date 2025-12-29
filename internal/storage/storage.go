package storage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/shopspring/decimal"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrInsufficientFunds  = errors.New("insufficient funds")
)

// Balance stor interface
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID string) (models.BalanceResponse, error)
	Withdraw(ctx context.Context, userID string, order string, sum float32) error
	GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error)
	AddAccrual(ctx context.Context, userID string, amount float32) error
}

// Order stor interface
type OrderRepository interface {
	GetOrderByNumber(ctx context.Context, userID string, orderNumber string) (models.Order, error)
	GetOrderByNumberAnyUser(ctx context.Context, orderNumber string) (models.Order, error)
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrdersByUser(ctx context.Context, userID string) ([]models.OrderResponse, error)
	GetOrdersForProcessing(ctx context.Context, limit int) ([]models.OrdersForProcessing, error)
	UpdateOrder(ctx context.Context, updateReq models.OrderUpdate) error
}

// User stor interface
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

// parseDecimal - helper function to parse sql.NullString to decimal.Decimal
// returns decimal.Zero if NullString is not valid
func parseDecimal(accrual sql.NullString) (decimal.Decimal, error) {
	if accrual.Valid {
		d, err := decimal.NewFromString(accrual.String)
		if err != nil {
			return decimal.Zero, err
		}
		return d, nil
	}
	return decimal.Zero, nil
}

// parseDecimal - helper function to parse sql.NullString to float32
// returns float32 if NullString is not valid
func parseFloat(accrual sql.NullString) (float32, error) {
	if accrual.Valid {
		f64, err := strconv.ParseFloat(accrual.String, 32)
		if err != nil {
			return 0.0, err
		}
		return float32(f64), nil
	}
	return 0.0, nil
}
