package main

import (
	"context"
	"log"
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
		log.Fatal("Failed to connect to database:", err)
	}
	log.Printf("Connect database at %s", cfg.DatabaseURI)
	defer store.Close()

	// Start HTTP server
	router := handlers.NewRouter(cfg, store)
	go func() {
		log.Printf("Starting server at %s", cfg.RunAdress)
		err = http.ListenAndServe(cfg.RunAdress, router)
		if err != nil {
			log.Fatal("Failed to start server: ", err)
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
	log.Println("Shutting down worker...")

	cancel()
	worker.Stop()

	log.Println("Worker shutdown complete")
}
