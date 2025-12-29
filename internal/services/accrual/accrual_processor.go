package accrual

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/storage"
)

type Processor struct {
	accrualClient *Client
	orderRepo     storage.OrderRepository
	balanceRepo   storage.BalanceRepository
	batchSize     int
	pollInterval  time.Duration
}

// NewProcessor creates a new Processor instance with the given accrual client and repositories.
// It sets default values for batch size and polling interval.
func NewProcessor(
	accrualClient *Client,
	orderRepo storage.OrderRepository,
	balanceRepo storage.BalanceRepository,
) *Processor {
	return &Processor{
		accrualClient: accrualClient,
		orderRepo:     orderRepo,
		balanceRepo:   balanceRepo,
		batchSize:     10,
		pollInterval:  5 * time.Second,
	}
}

// ProcessOrders continuously processes orders by polling the order repository at regular intervals.
// It stops processing when the provided context is canceled.
func (p *Processor) ProcessOrders(ctx context.Context) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Accrual processor stopped")
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

// processBatch retrieves a batch of orders and processes each one.
// Errors during processing are logged but do not stop the batch processing.
func (p *Processor) processBatch(ctx context.Context) {
	orders, err := p.orderRepo.GetOrdersForProcessing(ctx, p.batchSize)
	if err != nil {
		log.Printf("Error getting orders for processing: %v", err)
		return
	}
	for _, order := range orders {
		if err := p.processOrder(ctx, order); err != nil {
			log.Printf("Error processing order %s: %v", order.Number, err)
		}
	}
}

// processOrder retrieves accrual information for a single order and updates its status and user's balance accordingly.
// Errors are returned to the caller for logging.
// It uses the accrual client to get order info and updates the order and balance repositories.
func (p *Processor) processOrder(ctx context.Context, order models.OrdersForProcessing) error {
	accrualResp, err := p.accrualClient.GetOrderInfo(ctx, order.Number)
	if err != nil {
		return fmt.Errorf("failed to get order info: %w", err)
	}
	if accrualResp == nil {
		return nil
	}
	statusForOrder := models.MapAccrualToInternalStatus(accrualResp.Status)
	updateReq := models.OrderUpdate{
		Number:  order.Number,
		Status:  statusForOrder,
		Accrual: accrualResp.Accrual,
	}
	log.Printf("Got %s status, for user %s with accrual %f", statusForOrder, order.UserID, accrualResp.Accrual)
	if err := p.orderRepo.UpdateOrder(ctx, updateReq); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	if accrualResp.Status == models.MapAccrualToInternalStatus(models.AccrualStatusProcessed) && accrualResp.Accrual > 0 {
		log.Printf("Balance for user %s updated with %f", order.UserID, accrualResp.Accrual)
		if err := p.balanceRepo.AddAccrual(ctx, order.UserID, accrualResp.Accrual); err != nil {
			return fmt.Errorf("failed to add accrual to balance: %w", err)
		}
	}
	log.Printf("Order %s updated to status %s with accrual %f",
		order.Number, accrualResp.Status, accrualResp.Accrual)
	return nil
}
