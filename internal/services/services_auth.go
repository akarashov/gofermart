package services

import (
	"context"
	"errors"
	"regexp"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/pkg/hasher"
	"github.com/akarashov/gofermart/internal/storage"
)

type authService struct {
	userRepo storage.UserRepository
	hasher   hasher.PasswordHasher
}

func NewAuthService(userRepo storage.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher.NewBcryptHasher(),
	}
}

// checkPasswordComplexity checks if the password meets complexity requirements:
// at least 8 characters, at most 64 characters, contains at least one number,
// one uppercase letter, one lowercase letter, and one special character.
func checkPasswordComplexity(password string) bool {
	minLenght := len(password) >= 8
	maxLenght := len(password) <= 64
	hasNumber, _ := regexp.MatchString(`[0-9]`, password)
	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	hasLower, _ := regexp.MatchString(`[a-z]`, password)
	// hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
	return minLenght && maxLenght && hasNumber && hasUpper && hasLower
	// && hasSpecial
}

// checkLogin checks if the login meets requirements:
// only alphanumeric characters, underscores, dots, and dashes,
// length between 3 and 32 characters.
func checkLogin(login string) bool {
	valid, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]{3,32}$`, login)
	return valid
}

// Register registers a new user with the given login and password.
// It returns the user ID or an error if registration fails.
func (s *authService) Register(ctx context.Context, user models.UserRegister) (string, error) {
	if !checkLogin(user.Login) {
		return "", errors.New("login does not meet requirements: ^[a-zA-Z0-9_.-]{6,25}$")
	}
	if !checkPasswordComplexity(user.Password) {
		return "", errors.New("password does not meet complexity requirements")
	}
	hashedPassword, err := s.hasher.Hash(user.Password)
	if err != nil {
		return "", err
	}
	newUser := &models.User{
		Login:        user.Login,
		PasswordHash: hashedPassword,
	}
	return s.userRepo.CreateUser(ctx, newUser)
}

// Login authenticates a user with the given login and password.
// It returns the user ID or an error if authentication fails.
func (s *authService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", ErrWrongPasswordOrLogin
	}
	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return "", ErrWrongPasswordOrLogin
	}
	return user.ID, nil
}
