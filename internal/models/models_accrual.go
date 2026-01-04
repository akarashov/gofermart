package models

// Base response from accrual
type AccrualResponse struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

const (
	AccrualStatusRegistered = "REGISTERED"
	AccrualStatusProcessing = "PROCESSING"
	AccrualStatusProcessed  = "PROCESSED"
	AccrualStatusInvalid    = "INVALID"
)

// Accrual status to order status map
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
