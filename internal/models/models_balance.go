package models

// BalanceResponse - ответ для API при запросе баланса
type BalanceResponse struct {
	Current   float32 `json:"current"`   // Accrual sum
	Withdrawn float32 `json:"withdrawn"` // Withdrawn sum
}
