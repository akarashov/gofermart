package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *PostgresStorage) BalanceRepository() BalanceRepository {
	return p
}

// Get balance a specific user
// Sums up current and withdrawn from balance table for specific user
func (p *PostgresStorage) GetBalance(ctx context.Context, userID string) (models.BalanceResponse, error) {
	query := `SELECT current, withdrawn FROM balance WHERE user_id = $1`
	var balance models.BalanceResponse
	var currentStr sql.NullString
	var withdrawnStr sql.NullString
	err := p.db.QueryRowContext(ctx, query, userID).Scan(&currentStr, &withdrawnStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return balance, errors.New("balance not found")
		}
		return balance, fmt.Errorf("failed to get user by login: %w", err)
	}
	d, err := parseFloat(currentStr)
	if err != nil {
		return balance, fmt.Errorf("failed to parse current: %w", err)
	}
	balance.Current = d
	d, err = parseFloat(withdrawnStr)
	if err != nil {
		return balance, fmt.Errorf("failed to parse withdraw: %w", err)
	}
	balance.Withdrawn = d
	return balance, nil
}

// Withdraw - withdraw specific sum for specific user and order
// Uses withdraw stored procedure
func (p *PostgresStorage) Withdraw(ctx context.Context, userID string, order string, sum float32) error {
	query := `SELECT * FROM withdraw($1, $2, $3);`
	_, err := p.db.ExecContext(ctx, query, userID, order, sum)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23514" { // violates check constraint "check_current_non_negative"
			return ErrInsufficientFunds
		}
		return fmt.Errorf("failed to withdraw: %w", err)
	}
	return nil
}

// GetWithdrawals - get all withdrawals for specific user
// ordered by processed_at DESC
func (p *PostgresStorage) GetWithdrawals(ctx context.Context, userID string) ([]models.WithdrawalResponse, error) {
	query := `SELECT o.order_number, w.sum, w.processed_at FROM withdrawals w
	JOIN orders o ON w.order_id = o.id
	WHERE w.user_id = $1 ORDER BY w.processed_at DESC;`
	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawals: %w", err)
	}
	defer rows.Close()
	var withdrawals []models.WithdrawalResponse
	for rows.Next() {
		var wr models.WithdrawalResponse
		var withdrawnStr sql.NullString
		err := rows.Scan(&wr.Order, &withdrawnStr, &wr.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal: %w", err)
		}
		d, err := parseDecimal(withdrawnStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse withdrawn: %w", err)
		}
		wr.Sum = d
		withdrawals = append(withdrawals, wr)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return withdrawals, nil
}

// AddAccrual - add accrual amount to specific user's balance
// Assumes that the user already has a balance record
func (p *PostgresStorage) AddAccrual(ctx context.Context, userID string, amount float32) error {
	log.Printf("Adding accrual to DB of %f to user %s", amount, userID)
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE balance SET current=$1 WHERE user_id = $2`
	_, err = tx.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
