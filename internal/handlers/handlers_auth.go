package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/pkg/jwt"
	"github.com/akarashov/gofermart/internal/services"
	"github.com/akarashov/gofermart/internal/storage"
)

type AuthHandler struct {
	authService services.AuthService
	jwtManager  jwt.Manager
}

// NewAuthHandler creates a new AuthHandler with the given AuthService and JWT Manager.
// Returns a pointer to the created AuthHandler.
func NewAuthHandler(authService services.AuthService, jwtManager jwt.Manager) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtManager:  jwtManager,
	}
}

// Register handles user registration requests with JSON body containing login and password.
// Returns JWT token in Authorization header upon success.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userRegister models.UserRegister
	var buf bytes.Buffer
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &userRegister); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := h.authService.Register(r.Context(), userRegister)
	switch err {
	case nil:
		token, ok := h.jwtManager.Generate(userID)
		if ok != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		// w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", "BEARER "+token)
		w.WriteHeader(http.StatusOK)
		return
	case storage.ErrUserAlreadyExists:
		http.Error(w, err.Error(), http.StatusConflict)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Login handles user login requests with JSON body containing login and password.
// Returns JWT token in Authorization header upon success.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	var buf bytes.Buffer
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &userLogin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := h.authService.Login(r.Context(), userLogin.Login, userLogin.Password)
	switch err {
	case nil:
		token, ok := h.jwtManager.Generate(userID)
		if ok != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		w.Header().Set("Authorization", "BEARER "+token)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"jwt": token})
		return
	case services.ErrWrongPasswordOrLogin:
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
