package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// WithdrawalResponse - ответ при запросе истории списаний
type WithdrawalResponse struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
	ProcessedAt time.Time       `json:"processed_at"`
}

// WithdrawalRequest - Запрос на списание средств
type WithdrawalRequest struct {
	Order string `json:"order"`
	Sum   string `json:"sum"`
}
