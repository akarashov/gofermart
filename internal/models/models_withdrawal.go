package models

import (
	"time"
)

// Withdrawal Response for list request
type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// Withdrawal Request
type WithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
