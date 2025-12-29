package models

import (
	"time"
)

// WithdrawalResponse - ответ при запросе истории списаний
type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// WithdrawalRequest - Запрос на списание средств
type WithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
