package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/services"
	"github.com/go-chi/jwtauth"
)

// Base order struct
type OrdersHandler struct {
	orderService services.OrderService
}

// Create an instance of order
func NewOrdersHandler(orderService services.OrderService) *OrdersHandler {
	return &OrdersHandler{
		orderService: orderService,
	}
}

// UploadOrder handles the uploading of a new order.
// Expects Content-Type to be text/plain and the body to contain the order number.
// Returns appropriate HTTP status codes based on the outcome.
func (h *OrdersHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest) // - `400` — неверный формат запроса;
		return
	}
	token, claims, err := jwtauth.FromContext(r.Context())
	if token == nil || claims == nil || err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized) // - `401` — пользователь не аутентифицирован;
		return
	}
	userID := claims["sub"].(string)
	var orderUploadRequest models.OrderUploadRequest
	orderUploadRequestBinary, err := io.ReadAll(r.Body)
	orderUploadRequest.Number = string(orderUploadRequestBinary)
	if err != nil || orderUploadRequest.Number == "" {
		http.Error(w, "Error body parse", http.StatusBadRequest) // - `400` — неверный формат запроса;
	} else {
		err = h.orderService.UploadOrder(r.Context(), userID, orderUploadRequest.Number)
		switch err {
		case nil:
			w.WriteHeader(http.StatusAccepted) // - `202` — новый номер заказа принят в обработку;
			return
		case services.ErrOrderAlreadyUploadedByUser:
			w.WriteHeader(http.StatusOK) // - `200` — номер заказа уже был загружен этим пользователем;
			return
		case services.ErrOrderAlreadyUploadedByAnotherUser:
			http.Error(w, err.Error(), http.StatusConflict) // - `409` — номер заказа уже был загружен другим пользователем;
			return
		case services.ErrInvalidOrderNumber:
			http.Error(w, err.Error(), http.StatusUnprocessableEntity) // - `422` — неверный формат номера заказа;
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError) // - `500` — внутренняя ошибка сервера.
			return
		}
	}
}

// GetOrders retrieves all orders for the authenticated user.
// Returns appropriate HTTP status codes based on the outcome.
func (h *OrdersHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if token == nil || claims == nil || err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized) // - `401` — пользователь не аутентифицирован;
		return
	}
	userID := claims["sub"].(string)
	orders, err := h.orderService.GetUserOrders(r.Context(), userID)
	switch err {
	case nil:
		resp, err := json.Marshal(orders)
		if err != nil {
			http.Error(w, "Error on marshaling", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // - `200` — успешная обработка запроса.
		w.Write(resp)
		return
	case services.ErrNoOrdersForUser:
		http.Error(w, err.Error(), http.StatusNoContent) // - `204` — нет данных для ответа.
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError) // - `500` — внутренняя ошибка сервера.
		return
	}
}
