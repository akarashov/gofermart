package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage - create new PostgresStorage instance
// dataBaseDSN - database connection string
func NewPostgresStorage(dataBaseDSN string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dataBaseDSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	err = makeMigraton("file://migrations", dataBaseDSN)
	if err != nil {
		return &PostgresStorage{db: db}, err
	}
	return &PostgresStorage{db: db}, nil
}

// makeMigraton - runs database migrations
// using golang-migrate/migrate
// pathMigrations - path to migrations files
// dataBaseDSN - database connection string
func makeMigraton(pathMigrations string, dataBaseDSN string) error {
	m, err := migrate.New(pathMigrations, dataBaseDSN)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Ping - ping database to check connection
func (p *PostgresStorage) Ping(ctx context.Context) bool {
	err := p.db.PingContext(ctx)
	return err == nil
}

// Close - close database connection
func (p *PostgresStorage) Close() error {
	return p.db.Close()
}
