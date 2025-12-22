package handlers

import (
	"time"

	"github.com/akarashov/gofermart/internal/config"
	"github.com/akarashov/gofermart/internal/pkg/jwt"
	"github.com/akarashov/gofermart/internal/services"

	"github.com/akarashov/gofermart/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

// NewRouter - creates and configures a new chi.Mux router with all handlers and middleware
func NewRouter(cfg *config.Config, store *storage.PostgresStorage) *chi.Mux {
	authHandler := NewAuthHandler(
		services.NewAuthService(
			store.UserRepository()),
		jwt.NewJWTManager(
			cfg.JWTSecret,
			time.Duration(cfg.JWTExpireHours)*time.Hour))

	ordersHandler := NewOrdersHandler(
		services.NewOrderService(store.OrderRepository()))

	balanceHandler := NewBalanceHandler(
		services.NewBalanceService(store.BalanceRepository()))

	tokenAuth := jwtauth.New(
		"HS256",
		[]byte(cfg.JWTSecret),
		nil)

	Router := SetupRouter(
		authHandler,
		ordersHandler,
		balanceHandler,
		tokenAuth)
	return Router
}
