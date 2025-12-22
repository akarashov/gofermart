package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwtpkg "github.com/golang-jwt/jwt/v4"

	"github.com/akarashov/gofermart/internal/services"

	"github.com/akarashov/gofermart/internal/models"
	"github.com/akarashov/gofermart/internal/storage"
)

type mockAuthService struct {
	retErr error
	retID  string
	called bool
	got    models.UserRegister
}

func (m *mockAuthService) Register(ctx context.Context, user models.UserRegister) (string, error) {
	m.called = true
	m.got = user
	if m.retID != "" {
		return m.retID, m.retErr
	}
	return "123", m.retErr
}

func (m *mockAuthService) Login(ctx context.Context, login, password string) (string, error) {
	if m.retID != "" {
		return m.retID, m.retErr
	}
	return "123", m.retErr
}

type mockJWTManager struct {
	token  string
	err    error
	called bool
	gotID  string
}

func (m *mockJWTManager) Generate(userID string) (string, error) {
	m.called = true
	m.gotID = userID
	return m.token, m.err
}

func (m *mockJWTManager) Verify(tokenStr string) (*jwtpkg.RegisteredClaims, error) {
	return nil, nil
}

func TestAuthHandler_Register(t *testing.T) {
	goodUser := models.UserRegister{Login: "user", Password: "P@ssword1"}

	tests := []struct {
		name           string
		contentType    string
		body           interface{}
		authRetErr     error
		jwtToken       string
		jwtErr         error
		wantStatus     int
		wantHeader     string
		wantBodySubstr string
	}{
		{"success", "application/json", goodUser, nil, "token123", nil, http.StatusOK, "Bearer token123", ""},
		{"user_exists", "application/json", goodUser, storage.ErrUserAlreadyExists, "", nil, http.StatusConflict, "", "user already exists"},
		{"invalid_content_type", "text/plain", `not json`, nil, "", nil, http.StatusBadRequest, "", "Content-Type must be application/json"},
		{"invalid_json", "application/json", `{"`, nil, "", nil, http.StatusBadRequest, "", "unexpected end"},
		{"jwt_failure", "application/json", goodUser, nil, "", errors.New("jwt error"), http.StatusInternalServerError, "", "failed to generate token"},
		{"auth_error", "application/json", goodUser, errors.New("db fail"), "", nil, http.StatusInternalServerError, "", "db fail"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error
			switch v := tc.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(bodyBytes)))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rr := httptest.NewRecorder()

			mas := &mockAuthService{retErr: tc.authRetErr}
			mj := &mockJWTManager{token: tc.jwtToken, err: tc.jwtErr}
			h := NewAuthHandler(mas, mj)

			h.Register(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status: got %d want %d, body: %s", rr.Code, tc.wantStatus, rr.Body.String())
			}

			if tc.wantHeader != "" {
				got := rr.Header().Get("Authorization")
				if got != tc.wantHeader {
					t.Fatalf("Authorization header: got %q want %q", got, tc.wantHeader)
				}
			}

			if tc.wantBodySubstr != "" && !strings.Contains(rr.Body.String(), tc.wantBodySubstr) {
				t.Fatalf("body does not contain %q: %q", tc.wantBodySubstr, rr.Body.String())
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	goodLogin := models.UserLogin{Login: "user", Password: "P@ssword1"}

	tests := []struct {
		name           string
		contentType    string
		body           interface{}
		authRetErr     error
		jwtToken       string
		jwtErr         error
		wantStatus     int
		wantHeader     string
		wantBodySubstr string
	}{
		{"success", "application/json", goodLogin, nil, "token123", nil, http.StatusOK, "Bearer token123", ""},
		{"invalid_content_type", "text/plain", `not json`, nil, "", nil, http.StatusBadRequest, "", "Content-Type must be application/json"},
		{"invalid_json", "application/json", `{"`, nil, "", nil, http.StatusBadRequest, "", "unexpected end"},
		{"wrong_credentials", "application/json", goodLogin, services.ErrWrongPasswordOrLogin, "", nil, http.StatusUnauthorized, "", "wrong password or login"},
		{"auth_error", "application/json", goodLogin, errors.New("db fail"), "", nil, http.StatusInternalServerError, "", "db fail"},
		{"jwt_failure", "application/json", goodLogin, nil, "", errors.New("jwt error"), http.StatusInternalServerError, "", "failed to generate token"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			var err error
			switch v := tc.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(bodyBytes)))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rr := httptest.NewRecorder()

			mas := &mockAuthService{retErr: tc.authRetErr}
			mj := &mockJWTManager{token: tc.jwtToken, err: tc.jwtErr}
			h := NewAuthHandler(mas, mj)

			h.Login(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status: got %d want %d, body: %s", rr.Code, tc.wantStatus, rr.Body.String())
			}

			if tc.wantHeader != "" {
				got := rr.Header().Get("Authorization")
				if got != tc.wantHeader {
					t.Fatalf("Authorization header: got %q want %q", got, tc.wantHeader)
				}
			}

			if tc.wantBodySubstr != "" && !strings.Contains(rr.Body.String(), tc.wantBodySubstr) {
				t.Fatalf("body does not contain %q: %q", tc.wantBodySubstr, rr.Body.String())
			}
		})
	}
}
