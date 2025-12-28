package services

import (
	"context"
	"errors"

	"github.com/akarashov/gofermart/internal/models"
)

var (
	ErrWrongPasswordOrLogin              = errors.New("wrong password or login")
	ErrInvalidOrderNumber                = errors.New("invalid order number")
	ErrOrderAlreadyUploadedByUser        = errors.New("order already uploaded by this user")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrNoOrdersForUser                   = errors.New("no orders for this user")
	ErrInsufficientFunds                 = errors.New("insufficient funds")
	ErrNoWithdrawals                     = errors.New("no withdrawals found")
)

type AuthService interface {
	Register(ctx context.Context, user models.UserRegister) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
}

type BalanceService interface {
	GetBalance(ctx context.Context, userID string) (models.BalanceResponse, error)
	Withdraw(ctx context.Context, userID string, order string, sum float32) error
	GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error)
}

type OrderService interface {
	UploadOrder(ctx context.Context, userID, orderNumber string) error
	GetUserOrders(ctx context.Context, userID string) ([]models.OrderResponse, error)
}
