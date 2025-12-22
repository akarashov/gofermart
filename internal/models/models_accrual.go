package models

import "github.com/shopspring/decimal"

type AccrualResponse struct {
	Number  string          `json:"order"`             
	Status  string          `json:"status"`            
	Accrual decimal.Decimal `json:"accrual,omitempty"` 
}

const (
	AccrualStatusRegistered = "REGISTERED"
	AccrualStatusProcessing = "PROCESSING"
	AccrualStatusProcessed  = "PROCESSED"
	AccrualStatusInvalid    = "INVALID"
)

func MapAccrualToInternalStatus(accrualStatus string) string {
	switch accrualStatus {
	case AccrualStatusRegistered, AccrualStatusProcessing:
		return string(OrderStatusProcessing)
	case AccrualStatusProcessed:
		return string(OrderStatusProcessed)
	case AccrualStatusInvalid:
		return string(OrderStatusInvalid)
	default:
		return string(OrderStatusNew)
	}
}
