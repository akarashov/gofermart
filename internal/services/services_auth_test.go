package services

import (
	"context"
	"errors"
	"testing"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/storage"
)

// mockUserRepo is a small mock for storage.UserRepository used in tests.
type mockUserRepo struct {
	createFn func(ctx context.Context, user *models.User) (string, error)
	getFn    func(ctx context.Context, login string) (*models.User, error)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, user *models.User) (string, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	return "", nil
}

func (m *mockUserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	if m.getFn != nil {
		return m.getFn(ctx, login)
	}
	return nil, storage.ErrUserNotFound
}

func (m *mockUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return nil, storage.ErrUserNotFound
}

// fakeHasher implements hasher.PasswordHasher with predictable outputs for tests.
type fakeHasher struct{}

func (f *fakeHasher) Hash(password string) (string, error) {
	if password == "__err__" {
		return "", errors.New("hash failed")
	}
	return "hashed:" + password, nil
}

func (f *fakeHasher) Compare(hashedPassword, password string) error {
	if hashedPassword != "hashed:"+password {
		return errors.New("mismatch")
	}
	return nil
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		input     models.UserRegister
		repo      *mockUserRepo
		hasher    *fakeHasher
		wantID    string
		wantErr   bool
		wantErrIs error
	}{
		{
			name:  "success",
			input: models.UserRegister{Login: "user_ok", Password: "Aa1!aaaa"},
			repo: &mockUserRepo{createFn: func(ctx context.Context, user *models.User) (string, error) {
				if user.Login != "user_ok" {
					return "", errors.New("unexpected login")
				}
				if user.PasswordHash != "hashed:Aa1!aaaa" {
					return "", errors.New("unexpected hash")
				}
				return "u1", nil
			}},
			hasher:  &fakeHasher{},
			wantID:  "u1",
			wantErr: false,
		},
		{
			name:    "invalid login",
			input:   models.UserRegister{Login: "ab", Password: "Aa1!aaaa"},
			repo:    &mockUserRepo{},
			hasher:  &fakeHasher{},
			wantErr: true,
		},
		{
			name:    "invalid password",
			input:   models.UserRegister{Login: "valid_login", Password: "password"},
			repo:    &mockUserRepo{},
			hasher:  &fakeHasher{},
			wantErr: true,
		},
		{
			name:  "repo user exists",
			input: models.UserRegister{Login: "user_ok2", Password: "Aa1!aaaa"},
			repo: &mockUserRepo{createFn: func(ctx context.Context, user *models.User) (string, error) {
				return "", storage.ErrUserAlreadyExists
			}},
			hasher:    &fakeHasher{},
			wantErr:   true,
			wantErrIs: storage.ErrUserAlreadyExists,
		},
		{
			name:    "hasher error",
			input:   models.UserRegister{Login: "user_err", Password: "__err__"},
			repo:    &mockUserRepo{createFn: func(ctx context.Context, user *models.User) (string, error) { return "", nil }},
			hasher:  &fakeHasher{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{userRepo: tt.repo, hasher: tt.hasher}
			got, err := s.Register(ctx, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
					t.Fatalf("expected error to be %v, got %v", tt.wantErrIs, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantID {
				t.Fatalf("expected id %q, got %q", tt.wantID, got)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name    string
		login   string
		pass    string
		repo    *mockUserRepo
		hasher  *fakeHasher
		wantID  string
		wantErr error
	}{
		{
			name:  "success",
			login: "u_ok",
			pass:  "Secr3t!",
			repo: &mockUserRepo{getFn: func(ctx context.Context, login string) (*models.User, error) {
				return &models.User{ID: "u1", Login: login, PasswordHash: "hashed:Secr3t!"}, nil
			}},
			hasher: &fakeHasher{},
			wantID: "u1",
		},
		{
			name:  "wrong password",
			login: "u_ok",
			pass:  "badpass",
			repo: &mockUserRepo{getFn: func(ctx context.Context, login string) (*models.User, error) {
				return &models.User{ID: "u1", Login: login, PasswordHash: "hashed:Secr3t!"}, nil
			}},
			hasher:  &fakeHasher{},
			wantErr: ErrWrongPasswordOrLogin,
		},
		{
			name:    "user not found",
			login:   "no_user",
			pass:    "whatever",
			repo:    &mockUserRepo{getFn: func(ctx context.Context, login string) (*models.User, error) { return nil, storage.ErrUserNotFound }},
			hasher:  &fakeHasher{},
			wantErr: ErrWrongPasswordOrLogin,
		},
		{
			name:    "repo error",
			login:   "bad",
			pass:    "whatever",
			repo:    &mockUserRepo{getFn: func(ctx context.Context, login string) (*models.User, error) { return nil, errors.New("db down") }},
			hasher:  &fakeHasher{},
			wantErr: ErrWrongPasswordOrLogin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &authService{userRepo: tt.repo, hasher: tt.hasher}
			got, err := s.Login(ctx, tt.login, tt.pass)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.wantID {
				t.Fatalf("expected id %q, got %q", tt.wantID, got)
			}
		})
	}
}
