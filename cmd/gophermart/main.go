package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/akarashov/gofermart/internal/config"
	"github.com/akarashov/gofermart/internal/handlers"
	"github.com/akarashov/gofermart/internal/services/accrual"
	"github.com/akarashov/gofermart/internal/storage"
)

func main() {
	cfg := config.Load()
	store, err := storage.NewPostgresStorage(cfg.DatabaseURI)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer store.Close()

	// Start HTTP server
	router := handlers.NewRouter(cfg, store)
	go func() {
		slog.Info("server starting", "run_address", cfg.RunAddress)
		err = http.ListenAndServe(cfg.RunAddress, router)
		if err != nil {
			slog.Error("failed to start server", "err", err)
			os.Exit(1)
		}
	}()

	// Start accrual worker
	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress)

	worker := accrual.NewWorker(
		accrualClient,
		store.OrderRepository(),
		store.BalanceRepository(),
		cfg.Worker,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	slog.Info("worker shutting down")

	cancel()
	worker.Stop()

	slog.Info("worker shutdown complete")
}
