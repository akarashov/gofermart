package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

// SetupRouter configures the HTTP routes for the application.
// It sets up public and private routes with appropriate handlers and middleware.
func SetupRouter(
	authHandler *AuthHandler,
	ordersHandler *OrdersHandler,
	balanceHandler *BalanceHandler,
	tokenAuth *jwtauth.JWTAuth,
) *chi.Mux {
	mux := chi.NewRouter()
	mux.Group( //  Public routes
		func(mux chi.Router) {
			mux.Post("/api/user/register", authHandler.Register)
			mux.Post("/api/user/login", authHandler.Login)
		})
	mux.Group( //  Private routes
		func(mux chi.Router) {
			//  Authentication middleware
			mux.Use(jwtauth.Verifier(tokenAuth))
			mux.Use(jwtauth.Authenticator)
			//  Order routes
			mux.Post("/api/user/orders", ordersHandler.UploadOrder)
			mux.Get("/api/user/orders", ordersHandler.GetOrders)
			//  Balance routes
			mux.Get("/api/user/balance", balanceHandler.GetBalance)
			mux.Post("/api/user/balance/withdraw", balanceHandler.Withdraw)
			mux.Get("/api/user/withdrawals", balanceHandler.GetWithdrawals)
		})
	return mux
}
