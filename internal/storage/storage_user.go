package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *PostgresStorage) UserRepository() UserRepository {
	return p
}

// Create a new user
// Returns the new user's ID or an error
func (p *PostgresStorage) CreateUser(ctx context.Context, user *models.User) (string, error) {
	var userID string
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRowContext(ctx, query, user.Login, user.PasswordHash).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" { //  Unique constraint violation
			return "", ErrUserAlreadyExists
		}
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	balanceQuery := `INSERT INTO balance (user_id) VALUES ($1)`
	_, err = tx.ExecContext(ctx, balanceQuery, userID)
	if err != nil {
		return "", fmt.Errorf("failed to create user balance: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	return userID, nil
}

// Get user ID for Login
// Returns the user's ID or an error
func (p *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash FROM users WHERE login = $1`
	var user models.User
	err := p.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
// Returns the user or an error
func (p *PostgresStorage) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, login, password_hash FROM users WHERE id = $1`
	var user models.User
	err := p.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}
