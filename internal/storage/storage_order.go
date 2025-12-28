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

func (p *PostgresStorage) OrderRepository() OrderRepository {
	return p
}

// Get order by number for a specific user
// Returns the order or an error
func (p *PostgresStorage) GetOrderByNumber(ctx context.Context, userID string, orderNumber string) (models.Order, error) {
	query := `SELECT order_number, order_status, accrual, created_at, uploaded_at  FROM orders WHERE order_number = $1 and user_id = $2`
	var order models.Order
	var accrualStr sql.NullString
	err := p.db.QueryRowContext(
		ctx,
		query,
		orderNumber,
		userID).Scan(
		&order.Number,
		&order.Status,
		&accrualStr,
		&order.CreatedAt,
		&order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, errors.New("order not found")
		}
		return order, fmt.Errorf("failed to get order by number: %w", err)
	}
	d, err := parseDecimal(accrualStr)
	if err != nil {
		return order, fmt.Errorf("failed to parse accrual: %w", err)
	}
	order.Accrual = d
	return order, nil
}

// Get order by number for any user
// Returns the order or an error
func (p *PostgresStorage) GetOrderByNumberAnyUser(ctx context.Context, orderNumber string) (models.Order, error) {
	query := `SELECT order_number, order_status, accrual, created_at, uploaded_at FROM orders WHERE order_number = $1`
	var order models.Order
	var accrualStr sql.NullString
	err := p.db.QueryRowContext(
		ctx,
		query,
		orderNumber).Scan(
		&order.Number,
		&order.Status,
		&accrualStr,
		&order.CreatedAt,
		&order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, errors.New("order not found")
		}
		return order, fmt.Errorf("failed to get order by number: %w", err)
	}
	d, err := parseDecimal(accrualStr)
	if err != nil {
		return order, fmt.Errorf("failed to parse accrual: %w", err)
	}
	order.Accrual = d
	return order, nil
}

// Create a new order
// Returns an error if creation fails
func (p *PostgresStorage) CreateOrder(ctx context.Context, order models.Order) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `INSERT INTO orders (user_id, order_number) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, query, order.UserID, order.Number)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" { //  Unique constraint violation
			return ErrOrderAlreadyExists
		}
		return fmt.Errorf("failed to create order: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Get orders for a specific user
// Returns a list of orders or an error
func (p *PostgresStorage) GetOrdersByUser(ctx context.Context, userID string) ([]models.OrderResponse, error) {
	query := `SELECT order_number, order_status, accrual, uploaded_at  FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	var orders []models.OrderResponse
	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("orders not found")
		}
		return nil, fmt.Errorf("failed to get orders by user: %w", err)
	}
	for rows.Next() {
		var order models.OrderResponse
		var accrualStr sql.NullString
		err := rows.Scan(
			&order.Number,
			&order.Status,
			&accrualStr,
			&order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		d, err := parseFloat(accrualStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse accrual: %w", err)
		}
		order.Accrual = d
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return orders, nil
}

// Get orders for accrual processing in status NEW or PROCESSING by limit
func (p *PostgresStorage) GetOrdersForProcessing(ctx context.Context, limit int) ([]models.OrdersForProcessing, error) {
	// log.Printf("Getting up to %d orders for processing", limit)
	query := `SELECT user_id, order_number FROM orders WHERE order_status IN ('NEW', 'PROCESSING') ORDER BY uploaded_at ASC LIMIT $1`
	var orders []models.OrdersForProcessing
	rows, err := p.db.QueryContext(ctx, query, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("orders not found")
		}
		return nil, fmt.Errorf("failed to get orders by user: %w", err)
	}
	for rows.Next() {
		var order models.OrdersForProcessing
		err := rows.Scan(
			&order.UserID,
			&order.Number)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return orders, nil
}

// Update order from accrual
func (p *PostgresStorage) UpdateOrder(ctx context.Context, updateReq models.OrderUpdate) error {
	log.Printf("Updating order %s to status %s with accrual %f",
		updateReq.Number,
		updateReq.Status,
		updateReq.Accrual)
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `UPDATE orders SET order_status = $1, accrual = $2 WHERE order_number = $3`
	_, err = tx.ExecContext(ctx, query, updateReq.Status, updateReq.Accrual, updateReq.Number)
	if err != nil {
		return fmt.Errorf("failed to update orders: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
