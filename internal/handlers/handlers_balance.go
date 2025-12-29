package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/services"
	"github.com/go-chi/jwtauth"
)

// Master balance struct
type BalanceHandler struct {
	balanceService services.BalanceService
}

// Create an instance of balance
func NewBalanceHandler(balanceService services.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

// GetBalance retrieves the balance for the authenticated user.
// Returns appropriate HTTP status codes based on the outcome.
func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if token == nil || claims == nil || err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized) // - `401` — пользователь не аутентифицирован;
		return
	}
	userID := claims["sub"].(string)
	balance, err := h.balanceService.GetBalance(r.Context(), userID)
	switch err {
	case nil:
		resp, err := json.Marshal(balance)
		if err != nil {
			http.Error(w, "Error on marshaling", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // - `200` — успешная обработка запроса.
		w.Write(resp)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError) // - `500` — внутренняя ошибка сервера.
		return
	}
}

// GetBalance retrieves the balance for the authenticated user.
// Returns appropriate HTTP status codes based on the outcome.
func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if token == nil || claims == nil || err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized) // - `401` — пользователь не аутентифицирован;
		return
	}
	userID := claims["sub"].(string)
	var withdrawRequest models.WithdrawalRequest
	var buf bytes.Buffer
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &withdrawRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.balanceService.Withdraw(r.Context(), userID, withdrawRequest.Order, withdrawRequest.Sum)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK) // - `200` — успешная обработка запроса.
		return
	case services.ErrInvalidOrderNumber:
		http.Error(w, err.Error(), http.StatusUnprocessableEntity) // - `422` — неверный формат номера заказа;
		return
	case services.ErrInsufficientFunds:
		http.Error(w, err.Error(), http.StatusPaymentRequired) // - `402` — недостаточно средств на счету.
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError) // - `500` — внутренняя ошибка сервера.
		return
	}
}

// GetWithdrawals retrieves all withdrawals for the authenticated user.
// Returns appropriate HTTP status codes based on the outcome.
func (h *BalanceHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if token == nil || claims == nil || err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized) // - `401` — пользователь не аутентифицирован;
		return
	}
	userID := claims["sub"].(string)
	withdrawals, err := h.balanceService.GetWithdrawals(r.Context(), userID)
	switch err {
	case nil:
		resp, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(w, "Error on marshaling", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // - `200` — успешная обработка запроса.
		w.Write(resp)
		return
	case services.ErrNoWithdrawals:
		w.WriteHeader(http.StatusNoContent) // - `204` - нет ни одного списания.
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError) // - `500` — внутренняя ошибка сервера.
		return
	}
}
