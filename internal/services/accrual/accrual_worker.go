package accrual

import (
	"context"
	"log"
	"sync"

	"github.com/akarashov/gofermart/internal/storage"
)

type Worker struct {
	processors []*Processor
	wg         sync.WaitGroup
}

// NewWorker creates a new Worker instance with the specified number of processors.
// Each processor is initialized with the provided accrual client and repositories.
func NewWorker(
	accrualClient *Client,
	orderRepo storage.OrderRepository,
	balanceRepo storage.BalanceRepository,
	workerCount int,
) *Worker {
	worker := &Worker{}
	for i := 0; i < workerCount; i++ {
		processor := NewProcessor(accrualClient, orderRepo, balanceRepo)
		worker.processors = append(worker.processors, processor)
	}
	return worker
}

// Start initiates all accrual processors to begin processing orders.
// It spawns a separate goroutine for each processor and logs the start of the workers.
func (w *Worker) Start(ctx context.Context) {
	for _, processor := range w.processors {
		w.wg.Add(1)
		go func(p *Processor) {
			defer w.wg.Done()
			p.ProcessOrders(ctx)
		}(processor)
	}
	log.Printf("Started %d accrual workers", len(w.processors))
}

// Stop waits for all accrual processors to finish processing.
func (w *Worker) Stop() {
	w.wg.Wait()
	log.Println("All accrual workers stopped")
}
