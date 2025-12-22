package models

import (
	"github.com/shopspring/decimal"
)

// BalanceResponse - ответ для API при запросе баланса
type BalanceResponse struct {
	Current   decimal.Decimal `json:"current"`   // Accrual sum
	Withdrawn decimal.Decimal `json:"withdrawn"` // Withdrawn sum
}
