package services

import (
	"context"
	"fmt"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/pkg/validator"
	"github.com/akarashov/gofermart/internal/storage"
)

type balanceService struct {
	balanceRepo storage.BalanceRepository
}

func NewBalanceService(balanceRepo storage.BalanceRepository) BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
	}
}

// GetBalance - get balance a specific user
// Sums up current and withdrawn from balance table for specific user
func (s *balanceService) GetBalance(ctx context.Context, userID string) (models.BalanceResponse, error) {
	balance, err := s.balanceRepo.GetBalance(ctx, userID)
	if err != nil {
		return balance, fmt.Errorf("failed to get user balance: %w", err)
	}
	return balance, nil
}

// Withdraw - withdraw specific sum for specific user and order
// Validates order number using Luhn algorithm
// Returns ErrInsufficientFunds if not enough funds
func (s *balanceService) Withdraw(ctx context.Context, userID string, order string, sum string) error {
	if !validator.ValidateLuhn(order) {
		return ErrInvalidOrderNumber
	}
	err := s.balanceRepo.Withdraw(ctx, userID, order, sum)
	switch err {
	case nil:
		return nil
	case storage.ErrInsufficientFunds:
		return ErrInsufficientFunds
	default:
		return fmt.Errorf("failed to withdraw: %w", err)
	}
}

// GetWithdrawals - get all withdrawals for specific user
// ordered by processed_at DESC
// Returns ErrNoWithdrawals if no withdrawals found
func (s *balanceService) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
	withdrawals, err := s.balanceRepo.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals: %w", err)
	}
	if len(withdrawals) == 0 {
		return nil, ErrNoWithdrawals
	}
	return withdrawals, nil
}
