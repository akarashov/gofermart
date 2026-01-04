package models

// Response for a balance request
type BalanceResponse struct {
	Current   float32 `json:"current"`   // Accrual sum
	Withdrawn float32 `json:"withdrawn"` // Withdrawn sum
}
